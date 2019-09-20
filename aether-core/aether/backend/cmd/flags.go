package cmd

import (
	"aether-core/aether/services/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

// Instructions: If you're adding a flag anywhere in the app, please add it to flags struct, render flags, and flagschanged.

type flag struct {
	value   interface{}
	changed bool
}

// Struct for flags. When there's a new flag, add it here.
type flags struct {
	loggingLevel          flag // int
	orgName               flag // string
	appName               flag // string
	port                  flag // int
	externalIp            flag // string
	bootstrapIp           flag // string
	bootstrapPort         flag // int
	bootstrapType         flag // int
	syncAndQuit           flag // bool
	printToStdout         flag // bool
	metricsDebugMode      flag // bool
	swarmPlan             flag // string
	killTimeout           flag // int
	swarmNodeId           flag // int
	fpCheckEnabled        flag // bool
	sigCheckEnabled       flag // bool
	powCheckEnabled       flag // bool
	pageSigCheckEnabled   flag // bool
	tlsEnabled            flag // bool
	backendAPIPort        flag // int
	backendAPIPublic      flag // bool
	adminFeAddr           flag // string
	adminFePk             flag // string
	allowLocalhostRemotes flag // string
	imprint               flag // bool
	// Flags will be all lowercase in terminal input, heads up.
}

// When there's a new flag, add the parsing logic here underneath.
// I'm aware that this one sets the changed field, and yet, there's another method to check changed fields underneath that doesn't use this. It's because without using reflect package Go doesn't allow iteration over struct fields, and reflect, when used, does slow things down.
func renderFlags(cmd *cobra.Command) flags {
	var fl flags
	ll, err := cmd.Flags().GetInt("logginglevel")
	if err != nil && !strings.Contains(err.Error(), "flag accessed but not defined") {
		logging.LogCrash(err)
	}
	fl.loggingLevel.value = ll
	fl.loggingLevel.changed = cmd.Flags().Changed("logginglevel")
	on, err2 := cmd.Flags().GetString("orgname")
	if err2 != nil && !strings.Contains(
		err2.Error(), "flag accessed but not defined") {
		logging.LogCrash(err2)
	}
	fl.orgName.value = on
	fl.orgName.changed = cmd.Flags().Changed("orgname")
	an, err3 := cmd.Flags().GetString("appname")
	if err3 != nil && !strings.Contains(
		err3.Error(), "flag accessed but not defined") {
		logging.LogCrash(err3)
	}
	fl.appName.value = an
	fl.appName.changed = cmd.Flags().Changed("appname")

	p, err4 := cmd.Flags().GetInt("port")
	if err4 != nil && !strings.Contains(
		err4.Error(), "flag accessed but not defined") {
		logging.LogCrash(err4)
	}
	fl.port.value = p
	fl.port.changed = cmd.Flags().Changed("port")

	ei, err5 := cmd.Flags().GetString("externalip")
	if err5 != nil && !strings.Contains(
		err5.Error(), "flag accessed but not defined") {
		logging.LogCrash(err5)
	}
	fl.externalIp.value = ei
	fl.externalIp.changed = cmd.Flags().Changed("externalip")

	bi, err6 := cmd.Flags().GetString("bootstrapip")
	if err6 != nil && !strings.Contains(
		err6.Error(), "flag accessed but not defined") {
		logging.LogCrash(err6)
	}
	fl.bootstrapIp.value = bi
	fl.bootstrapIp.changed = cmd.Flags().Changed("bootstrapip")

	bp, err7 := cmd.Flags().GetInt("bootstrapport")
	if err7 != nil && !strings.Contains(
		err7.Error(), "flag accessed but not defined") {
		logging.LogCrash(err7)
	}
	fl.bootstrapPort.value = bp
	fl.bootstrapPort.changed = cmd.Flags().Changed("bootstrapport")

	bt, err8 := cmd.Flags().GetInt("bootstraptype")
	if err8 != nil && !strings.Contains(
		err8.Error(), "flag accessed but not defined") {
		logging.LogCrash(err8)
	}
	fl.bootstrapType.value = bt
	fl.bootstrapType.changed = cmd.Flags().Changed("bootstraptype")

	se, err9 := cmd.Flags().GetBool("syncandquit")
	if err9 != nil && !strings.Contains(
		err9.Error(), "flag accessed but not defined") {
		logging.LogCrash(err9)
	}
	fl.syncAndQuit.value = se
	fl.syncAndQuit.changed = cmd.Flags().Changed("syncandquit")

	so, err10 := cmd.Flags().GetBool("printtostdout")
	if err10 != nil && !strings.Contains(
		err10.Error(), "flag accessed but not defined") {
		logging.LogCrash(err10)
	}
	fl.printToStdout.value = so
	fl.printToStdout.changed = cmd.Flags().Changed("printtostdout")

	dm, err11 := cmd.Flags().GetBool("metricsdebugmode")
	if err11 != nil && !strings.Contains(
		err11.Error(), "flag accessed but not defined") {
		logging.LogCrash(err11)
	}
	fl.metricsDebugMode.value = dm
	fl.metricsDebugMode.changed = cmd.Flags().Changed("metricsdebugmode")

	sp, err12 := cmd.Flags().GetString("swarmplan")
	if err12 != nil && !strings.Contains(
		err12.Error(), "flag accessed but not defined") {
		logging.LogCrash(err12)
	}
	fl.swarmPlan.value = sp
	fl.swarmPlan.changed = cmd.Flags().Changed("swarmplan")

	kt, err13 := cmd.Flags().GetInt("killtimeout")
	if err13 != nil && !strings.Contains(
		err13.Error(), "flag accessed but not defined") {
		logging.LogCrash(err13)
	}
	fl.killTimeout.value = kt
	fl.killTimeout.changed = cmd.Flags().Changed("killtimeout")

	sni, err14 := cmd.Flags().GetInt("swarmnodeid")
	if err14 != nil && !strings.Contains(
		err14.Error(), "flag accessed but not defined") {
		logging.LogCrash(err14)
	}
	fl.swarmNodeId.value = sni
	fl.swarmNodeId.changed = cmd.Flags().Changed("swarmnodeid")

	fp, err15 := cmd.Flags().GetBool("fpcheckenabled")
	if err15 != nil && !strings.Contains(
		err15.Error(), "flag accessed but not defined") {
		logging.LogCrash(err15)
	}
	fl.fpCheckEnabled.value = fp
	fl.fpCheckEnabled.changed = cmd.Flags().Changed("fpcheckenabled")

	sig, err16 := cmd.Flags().GetBool("sigcheckenabled")
	if err16 != nil && !strings.Contains(
		err16.Error(), "flag accessed but not defined") {
		logging.LogCrash(err16)
	}
	fl.sigCheckEnabled.value = sig
	fl.sigCheckEnabled.changed = cmd.Flags().Changed("sigcheckenabled")

	pow, err17 := cmd.Flags().GetBool("powcheckenabled")
	if err17 != nil && !strings.Contains(
		err17.Error(), "flag accessed but not defined") {
		logging.LogCrash(err17)
	}
	fl.powCheckEnabled.value = pow
	fl.powCheckEnabled.changed = cmd.Flags().Changed("powcheckenabled")

	psig, err18 := cmd.Flags().GetBool("pagesigcheckenabled")
	if err18 != nil && !strings.Contains(
		err18.Error(), "flag accessed but not defined") {
		logging.LogCrash(err18)
	}
	fl.pageSigCheckEnabled.value = psig
	fl.pageSigCheckEnabled.changed = cmd.Flags().Changed("pagesigcheckenabled")

	tls, err19 := cmd.Flags().GetBool("tlsenabled")
	if err19 != nil && !strings.Contains(
		err19.Error(), "flag accessed but not defined") {
		logging.LogCrash(err19)
	}
	fl.tlsEnabled.value = tls
	fl.tlsEnabled.changed = cmd.Flags().Changed("tlsenabled")

	bport, err20 := cmd.Flags().GetInt("backendapiport")
	if err20 != nil && !strings.Contains(
		err20.Error(), "flag accessed but not defined") {
		logging.LogCrash(err20)
	}
	fl.backendAPIPort.value = bport
	fl.backendAPIPort.changed = cmd.Flags().Changed("backendapiport")

	bpub, err21 := cmd.Flags().GetBool("backendapipublic")
	if err21 != nil && !strings.Contains(
		err21.Error(), "flag accessed but not defined") {
		logging.LogCrash(err21)
	}
	fl.backendAPIPublic.value = bpub
	fl.backendAPIPublic.changed = cmd.Flags().Changed("backendapipublic")

	sfeaddr, err22 := cmd.Flags().GetString("adminfeaddr")
	if err22 != nil && !strings.Contains(
		err22.Error(), "flag accessed but not defined") {
		logging.LogCrash(err22)
	}
	fl.adminFeAddr.value = sfeaddr
	fl.adminFeAddr.changed = cmd.Flags().Changed("adminfeaddr")

	sfepk, err23 := cmd.Flags().GetString("adminfepk")
	if err23 != nil && !strings.Contains(
		err23.Error(), "flag accessed but not defined") {
		logging.LogCrash(err23)
	}
	fl.adminFePk.value = sfepk
	fl.adminFePk.changed = cmd.Flags().Changed("adminfepk")

	lhr, err24 := cmd.Flags().GetBool("allowlocalhostremotes")
	if err24 != nil && !strings.Contains(
		err24.Error(), "flag accessed but not defined") {
		logging.LogCrash(err24)
	}
	fl.allowLocalhostRemotes.value = lhr
	fl.allowLocalhostRemotes.changed = cmd.Flags().Changed("allowlocalhostremotes")

	impr, err25 := cmd.Flags().GetBool("imprint")
	if err25 != nil && !strings.Contains(
		err25.Error(), "flag accessed but not defined") {
		logging.LogCrash(err25)
	}
	fl.imprint.value = impr
	fl.imprint.changed = cmd.Flags().Changed("imprint")

	return fl
}

// When there's a new flag, add it underneath so that it'll be checked if a value was provided. If it is, we want to disable the writes.
func flagsChanged(cmd *cobra.Command) bool {
	var result bool
	// exceptions
	isProdFlag := func(name string) bool {
		return name == "logginglevel" ||
			name == "backendapiport" ||
			name == "backendapipublic" ||
			name == "adminfeaddr" ||
			name == "adminfepk"
	}
	changeChecker := func(flag *pflag.Flag) {
		if flag.Changed {
			if !isProdFlag(flag.Name) {
				// if not a prod flag, the app is running in debug mode, and we prevent writes to permanent config.
				result = true
			}
		}
	}
	cmd.Flags().VisitAll(changeChecker)
	return result

	// ... For every flag, we need this, because if a flag is given we need to stop writing to config store file, and only keep the config store object in memory.

	// What that means is that if you provide ANY flags, the app won't commit ANYTHING to the file - not just the flag you set, but anything else, too. It will effectively operate in read-only mode in terms of configuration. This read-only mode will activate only after the init of the configstore is complete, so it does not prevent initial creation or fixing of missing values.
	return false
}
