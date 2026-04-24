package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintStep(step int, total int, message string) {
	blue := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s %s\n", blue(fmt.Sprintf("[%d/%d]", step, total)), message)
}

func PrintSuccess(message string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("   %s %s\n", green("✔"), message)
}

func PrintError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("   %s %s\n", red("✘"), message)
}

func PrintInfo(message string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("   %s %s\n", yellow("ℹ"), message)
}

func Bold(message string) string {
	return color.New(color.Bold).Sprint(message)
}
