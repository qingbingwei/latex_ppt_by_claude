package latex

// GetTemplate returns a LaTeX Beamer template by name
func GetTemplate(name string) string {
	templates := map[string]string{
		"default": defaultTemplate,
		"madrid":  madridTemplate,
		"modern":  modernTemplate,
	}

	if template, ok := templates[name]; ok {
		return template
	}
	return defaultTemplate
}

const defaultTemplate = `\documentclass[aspectratio=169,11pt]{beamer}

% 中文支持
\usepackage[UTF8]{ctex}

% 主题设置
\usetheme{Madrid}
\usecolortheme{dolphin}

% 常用宏包
\usepackage{amsmath, amssymb, amsfonts}
\usepackage{graphicx}
\usepackage{hyperref}
\usepackage{booktabs}
\usepackage{ulem}

% 元信息
\title{{{TITLE}}}
\subtitle{{{SUBTITLE}}}
\author{{{AUTHOR}}}
\institute{{{INSTITUTE}}}
\date{\today}

\begin{document}

\begin{frame}
  \titlepage
\end{frame}

\begin{frame}{目录}
  \tableofcontents
\end{frame}

{{CONTENT}}

\end{document}`

const madridTemplate = `\documentclass[aspectratio=169,11pt]{beamer}

\usepackage[UTF8]{ctex}
\usetheme{Madrid}
\usecolortheme{whale}

\usepackage{amsmath, amssymb, amsfonts}
\usepackage{graphicx}
\usepackage{hyperref}
\usepackage{booktabs}

\title{{{TITLE}}}
\subtitle{{{SUBTITLE}}}
\author{{{AUTHOR}}}
\date{\today}

\begin{document}

\begin{frame}
  \titlepage
\end{frame}

\begin{frame}{目录}
  \tableofcontents
\end{frame}

{{CONTENT}}

\end{document}`

const modernTemplate = `\documentclass[aspectratio=169,11pt]{beamer}

\usepackage[UTF8]{ctex}
\usetheme{metropolis}

\usepackage{amsmath, amssymb, amsfonts}
\usepackage{graphicx}
\usepackage{hyperref}
\usepackage{booktabs}

\title{{{TITLE}}}
\subtitle{{{SUBTITLE}}}
\author{{{AUTHOR}}}
\date{\today}

\begin{document}

\begin{frame}
  \titlepage
\end{frame}

\begin{frame}{目录}
  \tableofcontents
\end{frame}

{{CONTENT}}

\end{document}`

func ListTemplates() []string {
	return []string{"default", "madrid", "modern"}
}
