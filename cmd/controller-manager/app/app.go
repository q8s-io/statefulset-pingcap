package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/util/term"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/leaderelection"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/version"
	"k8s.io/klog"
	utilflag "k8s.io/kubernetes/pkg/util/flag"

	pcinformers "github.com/q8s-io/statefulset-pingcap/client/client/informers/externalversions"
	"github.com/q8s-io/statefulset-pingcap/cmd/controller-manager/config"
	"github.com/q8s-io/statefulset-pingcap/cmd/controller-manager/options"
	"github.com/q8s-io/statefulset-pingcap/pkg/controller/statefulset"
)

// ResyncPeriod returns a function which generates a duration each time it is
// invoked; this is so that multiple controllers don't get into lock-step and all
// hammer the apiserver with list requests simultaneously.
func ResyncPeriod(c *config.CompletedConfig) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(c.GenericComponent.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run runs the controller-manager. This should never exit.
func Run(cc *config.CompletedConfig, stopCh <-chan struct{}) error {
	// To help debugging, immediately log version
	klog.Infof("Version: %+v", version.Get())

	run := func(ctx context.Context) {
		informerFactory := informers.NewSharedInformerFactory(cc.Client, cc.GenericComponent.MinResyncPeriod.Duration)
		pcInformerFactory := pcinformers.NewSharedInformerFactory(cc.PCClient, cc.GenericComponent.MinResyncPeriod.Duration)
		stsCtrl := statefulset.NewStatefulSetController(
			informerFactory.Core().V1().Pods(),
			pcInformerFactory.Apps().V1().StatefulSets(),
			informerFactory.Core().V1().PersistentVolumeClaims(),
			informerFactory.Apps().V1().ControllerRevisions(),
			cc.Client,
			cc.PCClient,
		)
		go stsCtrl.Run(runtime.NumCPU(), ctx.Done())
		// Start informers after all event listeners are registered.
		informerFactory.Start(ctx.Done())
		pcInformerFactory.Start(ctx.Done())
		<-ctx.Done()
	}

	ctx, cancel := context.WithCancel(context.TODO()) // TODO once Run() accepts a context, it should be used here
	defer cancel()

	go func() {
		select {
		case <-stopCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	// If leader election is enabled, runCommand via LeaderElector until done and exit.
	if cc.LeaderElection != nil {
		cc.LeaderElection.Callbacks = leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				klog.Fatalf("leaderelection lost")
			},
		}
		leaderElector, err := leaderelection.NewLeaderElector(*cc.LeaderElection)
		if err != nil {
			return fmt.Errorf("couldn't create leader elector: %v", err)
		}

		leaderElector.Run(ctx)

		return fmt.Errorf("lost lease")
	}

	run(ctx)
	return fmt.Errorf("finished without leader elect")
}

func NewControllerManagerCommand() *cobra.Command {
	opts := options.NewControllerManagerOptions()
	cmd := &cobra.Command{
		Use:  "controller-manager",
		Long: `Advanced StatefulSet Controller Manager`,
		Run: func(cmd *cobra.Command, args []string) {
			utilflag.PrintFlags(flag.CommandLine)

			c, err := opts.Config()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := Run(c.Complete(), wait.NeverStop); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	namedFlagSets := opts.Flags()
	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())
	for _, f := range namedFlagSets.FlagSets {
		flag.CommandLine.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})

	return cmd
}
