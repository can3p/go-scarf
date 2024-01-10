# go-scarf

A package to help you do scaffolding for you projects. Sometimes you have a
folder with many templates that you want to process to generate an example
project. This package allows you to do that with ease.

## Example usage

A file in your folder with templates

```
package templates // template folder contains all your templates

import "embed"

//go:embed *
var Template embed.FS
```

Scaffolding code:

```
s := scaffolder.New()

if !test {
	s = s.WithProcessor(scaffolder.FSProcessor(out))
}

return s.Scaffold(templates.Template, scaffolder.ScaffoldData{
	"ProjectName": projectName,
	"GithubRepo":  githubRepo,
})
```

### Features

* You can define your own processor, the default one would just write to stdout, `FProcessor` would create a folder and write all files there
* A custom funcmap can be specified, the default one only has `template_code` helper that prevents code escaping

### Conventions

* Only files ending with `.tmpl` are processed
* This suffix is dropped in the final result

## License

Apache 2.0
