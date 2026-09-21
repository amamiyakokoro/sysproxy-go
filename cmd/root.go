package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/amamiyakokoro/sysproxy-go/sysproxy"

	"github.com/spf13/cobra"
)

var (
	server     string
	bypass     string
	pacUrl     string
	waitServer bool

	device           string
	onlyActiveDevice bool
	multiThread      bool
	useRegistry      bool
	targetUser       string
	targetSID        string
	peerPID          int
	peerUID          uint32
	peerGID          uint32
	peerEnv          []string
)

const (
	waitServerDialTimeout = 500 * time.Millisecond
	waitServerInitialPoll = time.Second
	waitServerMaxPoll     = 30 * time.Second
	standaloneDeprecation = "DEPRECATED: the standalone sysproxy CLI is retained for compatibility only; KokoroBox integrations must use KokoroBox Service"
)

var rootCmd = &cobra.Command{
	Use:   "sysproxy",
	Short: "系统代理设置工具（已弃用；KokoroBox 请使用 KokoroBox Service）",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.ErrOrStderr(), standaloneDeprecation)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "设置系统代理",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		opts, err := commandOptions()
		if err != nil {
			fmt.Println("解析命令参数失败：", err)
			return
		}
		opts.Proxy = server
		opts.Bypass = bypass

		ctx := context.Background()
		if waitServer {
			var stop context.CancelFunc
			ctx, stop = signal.NotifyContext(ctx, os.Interrupt)
			defer stop()
		}

		err = setProxy(ctx, opts)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			fmt.Println("设置代理失败：", err)
			return
		}
		fmt.Println("代理设置成功，耗时：", time.Since(t))
	},
}

var pacCmd = &cobra.Command{
	Use:   "pac",
	Short: "设置 PAC 代理",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		opts, err := commandOptions()
		if err != nil {
			fmt.Println("解析命令参数失败：", err)
			return
		}
		opts.PACURL = pacUrl
		err = sysproxy.SetPac(opts)
		if err != nil {
			fmt.Println("设置 PAC 代理失败：", err)
			return
		}
		fmt.Println("PAC 代理设置成功，耗时：", time.Since(t))
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "取消代理设置",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		opts, err := commandOptions()
		if err != nil {
			fmt.Println("解析命令参数失败：", err)
			return
		}
		err = sysproxy.DisableProxy(opts)
		if err != nil {
			fmt.Println("取消代理设置失败：", err)
			return
		}
		fmt.Println("代理设置已取消，耗时：", time.Since(t))
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看当前代理设置",
	Run: func(cmd *cobra.Command, args []string) {
		opts, err := commandOptions()
		if err != nil {
			fmt.Println("解析命令参数失败：", err)
			return
		}
		status, err := sysproxy.QueryProxySettings(opts)
		if err != nil {
			fmt.Println("查询代理设置失败：", err)
			return
		}
		statusJSON, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			fmt.Println("格式化 JSON 失败：", err)
			return
		}
		fmt.Println(string(statusJSON))
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "监听系统代理设置变更",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		opts, err := commandOptions()
		if err != nil {
			fmt.Println("解析命令参数失败：", err)
			return
		}

		for {
			if err := sysproxy.WaitProxySettingsChange(ctx, opts); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Println("监听代理设置失败：", err)
				return
			}
			fmt.Println("update")
		}
	},
}

func setProxy(ctx context.Context, opts *sysproxy.Options) error {
	if waitServer {
		if err := waitProxyServerAvailable(ctx, opts.Proxy); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return sysproxy.SetProxy(opts)
}

func waitProxyServerAvailable(ctx context.Context, proxy string) error {
	addr, err := proxyTCPAddress(proxy)
	if err != nil {
		return err
	}
	fmt.Printf("等待代理服务器可用：%s\n", addr)

	dialer := net.Dialer{Timeout: waitServerDialTimeout}
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	nextPoll := time.Duration(0)
	for {
		if nextPoll > 0 {
			timer.Reset(nextPoll)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
			}
		}

		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			_ = conn.Close()
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if nextPoll == 0 {
			nextPoll = waitServerInitialPoll
		} else {
			nextPoll *= 2
			if nextPoll > waitServerMaxPoll {
				nextPoll = waitServerMaxPoll
			}
		}
	}
}

func proxyTCPAddress(proxy string) (string, error) {
	proxy = strings.TrimSpace(proxy)
	if proxy == "" {
		return "", fmt.Errorf("--wait-server 需要同时指定 --server")
	}
	if _, _, err := net.SplitHostPort(proxy); err == nil {
		return proxy, nil
	}
	return "", fmt.Errorf("invalid proxy address: %s", proxy)
}

func commandOptions() (*sysproxy.Options, error) {
	opts := &sysproxy.Options{
		Device:           device,
		OnlyActiveDevice: onlyActiveDevice,
		Concurrent:       sysproxy.Bool(multiThread),
		UseRegistry:      useRegistry,
	}
	if err := applyTargetOptions(opts); err != nil {
		return nil, err
	}
	return opts, nil
}

func applyTargetOptions(opts *sysproxy.Options) error {
	switch runtime.GOOS {
	case "windows":
		return applyWindowsTargetOptions(opts)
	case "linux":
		return applyLinuxTargetOptions(opts)
	default:
		return nil
	}
}

func applyWindowsTargetOptions(opts *sysproxy.Options) error {
	if targetUser != "" {
		userOpts, err := sysproxy.OptionsForUser(targetUser)
		if err != nil {
			return err
		}
		mergeTargetOptions(opts, userOpts)
	}
	if peerPID > 0 {
		processOpts, err := sysproxy.OptionsForProcess(peerPID)
		if err != nil {
			return err
		}
		mergeTargetOptions(opts, processOpts)
	}
	if targetSID != "" {
		opts.UserSID = targetSID
	}
	return nil
}

func applyLinuxTargetOptions(opts *sysproxy.Options) error {
	if targetUser != "" {
		userOpts, err := sysproxy.OptionsForUser(targetUser)
		if err != nil {
			return err
		}
		mergeTargetOptions(opts, userOpts)
	}
	if peerPID > 0 {
		processOpts, err := sysproxy.OptionsForProcess(peerPID)
		if err != nil {
			return err
		}
		mergeTargetOptions(opts, processOpts)
	}
	if peerUID != 0 {
		opts.PeerUID = peerUID
	}
	if peerGID != 0 {
		opts.PeerGID = peerGID
	}
	if len(peerEnv) > 0 {
		for _, item := range peerEnv {
			if !strings.Contains(item, "=") {
				return fmt.Errorf("--env 需要 KEY=VALUE：%s", item)
			}
		}
		opts.Environment = append([]string(nil), peerEnv...)
	}
	return nil
}

func mergeTargetOptions(dst, src *sysproxy.Options) {
	if dst == nil || src == nil {
		return
	}
	if src.UserSID != "" {
		dst.UserSID = src.UserSID
	}
	if src.PeerPID != 0 {
		dst.PeerPID = src.PeerPID
	}
	if src.PeerUID != 0 {
		dst.PeerUID = src.PeerUID
	}
	if src.PeerGID != 0 {
		dst.PeerGID = src.PeerGID
	}
	if len(src.Environment) > 0 {
		dst.Environment = append([]string(nil), src.Environment...)
	}
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(pacCmd)
	rootCmd.AddCommand(disableCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(watchCmd)

	rootCmd.PersistentFlags().BoolVarP(&onlyActiveDevice, "only-active-device", "a", false, "仅对活跃的网络设备生效")
	rootCmd.PersistentFlags().StringVarP(&device, "device", "d", "", "指定网络设备")
	rootCmd.PersistentFlags().BoolVar(&multiThread, "multithread", sysproxy.DefaultConcurrent(), "启用多线程并发设置；macOS 默认开启，Windows 默认关闭")
	rootCmd.PersistentFlags().BoolVar(&useRegistry, "registry", false, "Windows 使用注册表设置/查询代理，不调用 win32 API")
	registerTargetFlags()

	proxyCmd.Flags().StringVarP(&server, "server", "s", "", "代理服务器地址")
	proxyCmd.Flags().StringVarP(&bypass, "bypass", "b", "", "绕过地址")
	proxyCmd.Flags().BoolVar(&waitServer, "wait-server", false, "设置系统代理前一直等待代理服务器可用")

	pacCmd.Flags().StringVarP(&pacUrl, "url", "u", "", "pac 地址")

}

func registerTargetFlags() {
	switch runtime.GOOS {
	case "windows":
		rootCmd.PersistentFlags().StringVar(&targetUser, "user", "", "指定系统用户")
		rootCmd.PersistentFlags().StringVar(&targetSID, "sid", "", "Windows 指定用户 SID")
		rootCmd.PersistentFlags().IntVar(&peerPID, "pid", 0, "指定会话进程 PID")
	case "linux":
		rootCmd.PersistentFlags().StringVar(&targetUser, "user", "", "指定系统用户")
		rootCmd.PersistentFlags().IntVar(&peerPID, "pid", 0, "指定会话进程 PID")
		rootCmd.PersistentFlags().Uint32Var(&peerUID, "uid", 0, "Linux 指定用户 UID")
		rootCmd.PersistentFlags().Uint32Var(&peerGID, "gid", 0, "Linux 指定用户 GID")
		rootCmd.PersistentFlags().StringArrayVar(&peerEnv, "env", nil, "Linux 指定会话环境变量 KEY=VALUE，可重复")
	}
}
