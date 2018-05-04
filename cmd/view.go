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
	"github.com/jlucasnsilva/atog/ui"
	"github.com/spf13/cobra"
)

var (
	uiParams ui.Params
)

var viewCmd = &cobra.Command{
	Use:     "view filename*",
	Short:   "Watch the log files.",
	Long:    `Pass to this command a list of files to be watched.`,
	Example: "atog view example.log",
	Run: func(cmd *cobra.Command, args []string) {
		ps := uiParams
		ps.Filenames = args
		ui.Execute(ps)
	},
}

func init() {
	flags := viewCmd.Flags()

	flags.UintVarP(&uiParams.BufferSize, "buffer-size", "s", 10, "Maximum number of log messages to be displayed.")
	flags.BoolVarP(&uiParams.Empty, "empty", "e", false, "If true, no log messages prior to the openning of atog will\nbe displayed.")
}
