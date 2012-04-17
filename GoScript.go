// Copyright 2012  The "GoScript" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	// Sadly, local packages aren't legal anymore (see http://codereview.appspot.com/5787055)
	// so we can't just use "./gojs"; suggested workarounds apparently involve
	// mimicking the original source directory... or having all forkers hand-edit
	// their code to get it to compile. Bleah.
	gojs "github.com/steven-johnson/GoScript/gojs"
	"os"
)

func main() {
	var srcFile *string = flag.String("src", "", "Source js file to compile into go")
	flag.Parse()
	if len(*srcFile) == 0 {
		fmt.Println("Must specify srcFile")
		flag.PrintDefaults()
		os.Exit(1)
	}
	gojs.Compile(*srcFile)
}
