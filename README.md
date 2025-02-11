A collection of utilities for parsing go templates from directories,
with support for template inheritance.

Given a directory structure like
```
super.tmpl
home.tmpl
products.tmpl
admin/
  super.tmpl
  admin.tmpl
```

It will create the templates
- `super.tmpl`
- `home.tmpl`
- `products.tmpl`
- `admin/super.tmpl`
- `admin/admin.tmpl`
  
where `home.tmpl`, `products.tmpl` and `admin/admin.tmpl` are associated with a copy of `super.tmpl`,
and `admin/admin.tmpl` is associated with a copy of `admin/super.tmpl`.

This can be used to define common template fragments on the `super.tmpl` files, or to simulate a kind of template inheritance.

For example, by defining `super.tmpl` as
```
<html>
  <head>
    <title>Website</title>
  </head>
  <body>
    <main>
      {{ block "content" .}}{{end}}
    </main>
  </body>
</html>
```
templates lower in the chain can define `content` and include `super.tmpl` to reuse the page definition:
```
-- home.tmpl --
{{template "super.tmpl" .}}
{{define "content"}}
This is the content of home
{{end}}
```
renders as
```
<html>
  <head>
    <title>Website</title>
  </head>
  <body>
    <main>
      This is the content of home
    </main>
  </body>
</html>
```

To learn more about how to write Go templates, see https://pkg.go.dev/html/template. 
