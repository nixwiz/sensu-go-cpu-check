package main

import (
	"fmt"
	"math"
	"time"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
	"github.com/shirou/gopsutil/cpu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Warning  int
	Critical int
	Interval int64
}

type CPUUsage struct {
	Idle    float64
	System  float64
	User    float64
	Nice    float64
	Iowait  float64
	Irq     float64
	Softirq float64
	Steal   float64
	Cores   int
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-cpu-check",
			Short:    "Sensu Go CPU Check",
			Keyspace: "sensu.io/plugins/sensu-go-cpu-check/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		&sensu.PluginConfigOption{
			Path:      "warning",
			Env:       "",
			Argument:  "warning",
			Shorthand: "w",
			Default:   80,
			Usage:     "Warning threshold (>=) for CPU usage",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption{
			Path:      "critical",
			Env:       "",
			Argument:  "critical",
			Shorthand: "c",
			Default:   90,
			Usage:     "Critical threshold (>=) for CPU usage",
			Value:     &plugin.Critical,
		},
		&sensu.PluginConfigOption{
			Path:      "interval",
			Env:       "",
			Argument:  "interval",
			Shorthand: "i",
			Default:   int64(1),
			Usage:     "How long to sleep between CPU usage samples, in seconds",
			Value:     &plugin.Interval,
		},
	}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if plugin.Warning > plugin.Critical {
		return sensu.CheckStateWarning, fmt.Errorf("warning argument cannot be greater than critical argument")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	interval := time.Second * time.Duration(plugin.Interval)
	cpuTimes1, _ := cpu.Times(false)
	time.Sleep(interval)
	cpuTimes2, _ := cpu.Times(false)
	cpuTimes1All := getAll(cpuTimes1[0])
	cpuTimes2All := getAll(cpuTimes2[0])

	usage := new(CPUUsage)

	usage.Idle = getUsage(cpuTimes2[0].Idle, cpuTimes1[0].Idle, cpuTimes2All, cpuTimes1All)
	usage.System = getUsage(cpuTimes2[0].System, cpuTimes1[0].System, cpuTimes2All, cpuTimes1All)
	usage.User = getUsage(cpuTimes2[0].User, cpuTimes1[0].User, cpuTimes2All, cpuTimes1All)
	usage.Nice = getUsage(cpuTimes2[0].Nice, cpuTimes1[0].Nice, cpuTimes2All, cpuTimes1All)
	usage.Iowait = getUsage(cpuTimes2[0].Iowait, cpuTimes1[0].Iowait, cpuTimes2All, cpuTimes1All)
	usage.Irq = getUsage(cpuTimes2[0].Irq, cpuTimes1[0].Irq, cpuTimes2All, cpuTimes1All)
	usage.Softirq = getUsage(cpuTimes2[0].Softirq, cpuTimes1[0].Softirq, cpuTimes2All, cpuTimes1All)
	usage.Steal = getUsage(cpuTimes2[0].Steal, cpuTimes1[0].Steal, cpuTimes2All, cpuTimes1All)

	// Cheap way to get the core count
	cpuCores, _ := cpu.Times(true)
	usage.Cores = len(cpuCores)

	busy := usage.System + usage.User + usage.Nice + usage.Iowait + usage.Irq + usage.Softirq + usage.Steal

	if busy >= float64(plugin.Critical) {
		fmt.Printf("%s CRITICAL: ", plugin.PluginConfig.Name)
		formatOutput(usage)
		return sensu.CheckStateCritical, nil
	} else if busy >= float64(plugin.Warning) {
		fmt.Printf("%s WARNING: ", plugin.PluginConfig.Name)
		formatOutput(usage)
		return sensu.CheckStateWarning, nil
	}
	fmt.Printf("%s OK: ", plugin.PluginConfig.Name)
	formatOutput(usage)
	return sensu.CheckStateOK, nil
}

func getAll(t cpu.TimesStat) float64 {
	return t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal + t.Idle
}

func getUsage(t1, t2, t1all, t2all float64) float64 {
	return math.Min(100, math.Max(0, (t2-t1)/(t2all-t1all)*100))
}

func formatOutput(u *CPUUsage) {
	fmt.Printf("idle=%.2f%% user=%.2f%% system=%.2f%% iowait=%.2f%% nice=%.2f%% ", u.Idle, u.User, u.System, u.Iowait, u.Nice)
	fmt.Printf("irq=%.2f%% softirq=%.2f%% steal=%.2f%%", u.Irq, u.Softirq, u.Steal)
	fmt.Printf(" | ")
	fmt.Printf("cpu_idle=%.2f cpu_user=%.2f cpu_system=%.2f cpu_iowait=%.2f cpu_nice=%.2f ", u.Idle, u.User, u.System, u.Iowait, u.Nice)
	fmt.Printf("cpu_irq=%.2f cpu_softirq=%.2f cpu_steal=%.2f cpu_cores=%d\n", u.Irq, u.Softirq, u.Steal, u.Cores)
}
