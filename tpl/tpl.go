package tpl

import heroscrape "github.com/v-braun/hero-scrape"

var htmlTpl = `
<!DOCTYPE html>
<html lang="en" class=" is-copy-enabled">
<head prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb# object: http://ogp.me/ns/object# article: http://ogp.me/ns/article# profile: http://ogp.me/ns/profile#">
  <meta charset='utf-8'>
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta http-equiv="Content-Language" content="en">
  <meta name="viewport" content="width=device-width">

  <title>{{.Description}}</title>
  <link rel="fluid-icon" href="https://github.com/fluidicon.png" title="GitHub">
  <meta content="{{.Image}}" name="twitter:image:src" />
  <meta content="@github" name="twitter:site" />
  <meta content="summary" name="twitter:card" />
  <meta content="{{.User}}/{{.Repo}}" name="twitter:title" />
  <meta content="{{.Description}}" name="twitter:description" />
  <meta content="{{.Image}}" property="og:image" />
  <meta content="GitHub" property="og:site_name" />
  <meta content="object" property="og:type" />
  <meta content="{{.User}}/{{.Repo}}" property="og:title" />
  <meta content="{{.ResolverUrl}}" property="og:url" />
  <meta content="{{.Description}}" property="og:description" />
  <link rel="icon" type="image/x-icon" href="https://assets-cdn.github.com/favicon.ico">
  <meta name="description" content="{{.Description}}">
</head>
<body>
  <script>
    document.location.href = "{{.Url}}";
  </script>
</body>
</html>
`

type Model struct {
	Image       string
	Url         string
	Description string
	User        string
	Repo        string
	ResolverUrl string
}

func GetTemplate() string {
	return htmlTpl
}

func NewModel(user string, repo string, res *heroscrape.SearchResult) *Model {
	m := new(Model)
	m.Image = res.Image
	m.Description = res.Description
	m.Repo = repo
	m.User = user
	m.Url = "https://github.com/" + user + "/" + repo
	m.ResolverUrl = "https://share-my-repo.vibr.app/" + user + "/" + repo

	return m
}
