/*
 * Copyright (c) 2018, João Lucas Nunes e Silva
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL JOÃO LUCAS NUNES E SILVA BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/jlucasnsilva/atog/tabbed"
	"github.com/jlucasnsilva/atog/view"
	"github.com/jlucasnsilva/atog/watch"
	"github.com/spf13/cobra"
)

var (
	splitView bool
)

var rootCmd = &cobra.Command{
	Use:   "atog",
	Short: "Displays highlighted text from stdin or a file",
	Long: `
Displays highlighted text from stdin when no argument is passed or displays
highlighted text from text file when one argument is passed.`,
	Example: `    atog main.go
	
    tail -f  example.log | atog`,
	Run: func(cmd *cobra.Command, args []string) {
		nArgs := len(args)
		switch nArgs {
		case 0:
			if splitView {
				tabbed.Show(os.Stdin)
			} else {
				watch.Show(os.Stdin)
			}
		case 1:
			file, err := os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			view.Show(file)
		default:
			fmt.Printf("Pass only one file at a time.")
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&splitView, "split", "s", false, "")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
