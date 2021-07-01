package main

import _ "embed"

//go:embed .templates/agpl-3.0.tmpl
var Agpl30Template string

//go:embed .templates/apache-2.0.tmpl
var Apache20Template string

//go:embed .templates/bsd-2-clause.tmpl
var Bsd2ClauseTemplate string

//go:embed .templates/bsd-3-clause.tmpl
var Bsd3ClauseTemplate string

//go:embed .templates/cc0-1.0.tmpl
var Cc010Template string

//go:embed .templates/epl-2.0.tmpl
var Epl20Template string

//go:embed .templates/free-art-1.3.tmpl
var FreeArt13Template string

//go:embed .templates/gpl-2.0.tmpl
var Gpl20Template string

//go:embed .templates/gpl-3.0.tmpl
var Gpl30Template string

//go:embed .templates/lgpl-2.1.tmpl
var Lgpl21Template string

//go:embed .templates/lgpl-3.0.tmpl
var Lgpl30Template string

//go:embed .templates/mit.tmpl
var MitTemplate string

//go:embed .templates/mpl-2.0.tmpl
var Mpl20Template string

//go:embed .templates/unlicense.tmpl
var UnlicenseTemplate string

//go:embed .templates/wtfpl.tmpl
var WtfplTemplate string
