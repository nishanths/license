package main

import _ "embed"

// Agpl30Template is the template for the Affero General Public License.
//go:embed .templates/agpl-3.0.tmpl
var Agpl30Template string

// Apache20Template is the template for the Apache License.
//go:embed .templates/apache-2.0.tmpl
var Apache20Template string

// Bsd2ClauseTemplate is the template for the 2-Clause BSD License.
//go:embed .templates/bsd-2-clause.tmpl
var Bsd2ClauseTemplate string

// Bsd3ClauseTemplate is the template for the 3-Clause BSD License.
//go:embed .templates/bsd-3-clause.tmpl
var Bsd3ClauseTemplate string

// Cc010Template is the template for the CC0 Public Domain Waiver.
//go:embed .templates/cc0-1.0.tmpl
var Cc010Template string

// Epl20Template is the template for the Eclipse Public License.
//go:embed .templates/epl-2.0.tmpl
var Epl20Template string

// FreeArt13Template is the template for the Free Art License.
//go:embed .templates/free-art-1.3.tmpl
var FreeArt13Template string

// Gpl20Template is the template for the GNU General Public License (V2).
//go:embed .templates/gpl-2.0.tmpl
var Gpl20Template string

// Gpl30Template is the template for the GNU General Public License (V3).
//go:embed .templates/gpl-3.0.tmpl
var Gpl30Template string

// Lgpl21Template is the template for the GNU Lesser General Public License (V2).
//go:embed .templates/lgpl-2.1.tmpl
var Lgpl21Template string

// Lgpl30Template is the template for the GNU Lesser General Public License (V3).
//go:embed .templates/lgpl-3.0.tmpl
var Lgpl30Template string

// MitTemplate is the template for the MIT / X11 License.
//go:embed .templates/mit.tmpl
var MitTemplate string

// Mpl20Template is the template for the Mozilla Public License.
//go:embed .templates/mpl-2.0.tmpl
var Mpl20Template string

// UnlicenseTemplate is the template for the Unlicense.
//go:embed .templates/unlicense.tmpl
var UnlicenseTemplate string

// WtfplTemplate is the template for the Do What The Fuck You Want To Public License.
//go:embed .templates/wtfpl.tmpl
var WtfplTemplate string
