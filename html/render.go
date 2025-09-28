package html

import (
	"zonetree/cache"
)

const HEAD = `
<html>
<head>
        <title>TREE</title>
        <meta charset="utf-8"> <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="example2.css">
</head>
<body>

<figure> 
  <figcaption>Example</figcaption>
  <ul class="tree">

`

const FOOT = `

  </ul>
</figure>

</body>
</html	>
`

const NODESTART = `
    <li>
	<span>
`

const NODEEND = `
	</span>
    </li>
`

type Node struct {
	Name     string `json:"Name"`
	Parent   *Node  `json:"-"`
	Children []Node `json:"Children"`
}

func Tabs(nr int) string {
	var tab string
	for i := 0; i <= nr; i++ {
		tab += "\t"
	}
	return tab
}

func DrawNode(z cache.Zone) string {
	var node string
	node = NODESTART

	node += NODEEND
	return node
}
