package grub

import (
	"strconv"
	"strings"
)

// Config details the configuration for grub
type Config struct {
	Name          string   `json:"name,omitempty"`
	Kernel        string   `json:"kernel"`
	Initramfs     string   `json:"initramfs,omitempty"`
	KernelArgs    string   `json:"kernel_args,omitempty"`
	Multiboot     string   `json:"multiboot_kernel,omitempty"`
	MultibootArgs string   `json:"multiboot_args,omitempty"`
	Modules       []string `json:"multiboot_modules,omitempty"`
}

// GetDefaultConfig - will find the grub configuration that will be booted by default
func GetDefaultConfig(grubcfg string) (cfg *Config) {
	configs, index := ParseGrubCfg(grubcfg)
	if configs == nil {
		return nil
	}
	return &configs[index]
}

// ParseGrubCfg will do a line-by-line examination of the grub config
func ParseGrubCfg(grubcfg string) (configs []Config, defaultConfig int64) {
	var err error
	inMenuEntry := false
	var cfg *Config
	for _, line := range strings.Split(grubcfg, "\n") {
		// remove all leading spaces as they are not relevant for the config
		// line
		line = strings.TrimLeft(line, " ")
		sline := strings.Fields(line)
		if len(sline) == 0 {
			continue
		}

		if sline[0] == "set" && strings.Contains(sline[1], "default") {
			// Find the default configuration
			defaultCfgEntry := strings.Split(sline[1], "=")
			if defaultCfgEntry[0] == "default" {
				defaultConfig, err = strconv.ParseInt(trimQuote(defaultCfgEntry[1]), 8, 0)
				if err != nil {
					defaultConfig = 0
				}
			}
		}
		if sline[0] == "menuentry" {
			// if a "menuentry", start a new boot config
			if cfg != nil {
				// save the previous boot config, if any
				if cfg.Kernel != "" && cfg.Initramfs != "" {
					// only consider valid boot configs, i.e. the ones that have
					// both kernel and initramfs
					configs = append(configs, *cfg)
				}
			}
			inMenuEntry = true
			cfg = new(Config)
			name := ""
			if len(sline) > 1 {
				name = strings.Join(sline[1:], " ")
				name = strings.Replace(name, `\$`, "$", -1)
				name = strings.Split(name, "--")[0]
			}
			cfg.Name = name
		} else if inMenuEntry {
			// otherwise look for kernel and initramfs configuration
			if len(sline) < 2 {
				// surely not a valid linux or initrd directive, skip it
				continue
			}
			if sline[0] == "linux" || sline[0] == "linux16" || sline[0] == "linuxefi" {
				kernel := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = strings.Replace(cmdline, `\$`, "$", -1)
				cfg.Kernel = kernel
				cfg.KernelArgs = cmdline
			} else if sline[0] == "initrd" || sline[0] == "initrd16" || sline[0] == "initrdefi" {
				initrd := sline[1]
				cfg.Initramfs = initrd
			} else if sline[0] == "multiboot" || sline[0] == "multiboot2" {
				multiboot := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = strings.Replace(cmdline, `\$`, "$", -1)
				cfg.Multiboot = multiboot
				cfg.MultibootArgs = cmdline
			} else if sline[0] == "module" || sline[0] == "module2" {
				module := sline[1]
				cmdline := strings.Join(sline[2:], " ")
				cmdline = strings.Replace(cmdline, `\$`, "$", -1)
				if cmdline != "" {
					module = module + " " + cmdline
				}
				cfg.Modules = append(cfg.Modules, module)
			}
		}
	}

	// append last kernel config if it wasn't already
	if inMenuEntry && cfg.Kernel != "" && cfg.Initramfs != "" {
		configs = append(configs, *cfg)
	}
	return
}

type grubVersion int

var (
	grubV1 grubVersion = 1
	grubV2 grubVersion = 2
)

func unquote(ver grubVersion, text string) string {
	if ver == grubV2 {
		// if grub2, unquote the string, as directives could be quoted
		// https://www.gnu.org/software/grub/manual/grub/grub.html#Quoting
		// TODO unquote everything, not just \$
		return strings.Replace(text, `\$`, "$", -1)
	}
	// otherwise return the unmodified string
	return text
}

func trimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}
