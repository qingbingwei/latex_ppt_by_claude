# 使用 LaTeX（Beamer）制作 PPT 的 Claude 工作流说明

本文档用于约定如何让 Claude 帮助我们用 LaTeX 制作演示文稿（PPT），主要基于 `beamer` 文档类。

---

## 1. 基本要求

1. 统一使用 `beamer` 文档类生成 PDF 幻灯片。
2. 所有 LaTeX 文稿必须是**完整可编译文档**（含导言区与 `\begin{document}...\end{document}`）。
3. 代码块必须标注语言为 `latex`，方便复制：
   ```latex
   % example here
   ```
4. 所有中文文档默认使用 `ctex` 宏包或 `xelatex` 编译环境（根据模板指定）。

---

## 2. 推荐 Beamer 模板

### 2.1 基础中文模板（ctex + beamer）

未指定模板时，默认使用下列简洁模板：

```latex
\documentclass[aspectratio=169,11pt]{beamer}

% 中文支持（基于 pdflatex + ctex）
\usepackage[UTF8]{ctex}

% 主题设置（可按需修改）
\usetheme{Madrid}
\usecolortheme{dolphin}

% 常用宏包
\usepackage{amsmath, amssymb, amsfonts}
\usepackage{graphicx}
\usepackage{hyperref}
\usepackage{booktabs}
\usepackage{ulem} % 删除线等

% 元信息
\title{演示文稿标题}
\subtitle{副标题（可选）}
\author{作者姓名}
\institute{单位或学校}
\date{\today}

\begin{document}

% 标题页
\begin{frame}
  \titlepage
\end{frame}

% 目录页（可选）
\begin{frame}{目录}
  \tableofcontents
\end{frame}

% 示例内容页
\section{第一部分}
\begin{frame}{本页标题}
  \begin{itemize}
    \item 要点一
    \item 要点二
    \item 要点三
  \end{itemize}
\end{frame}

\end{document}
```

在此模板基础上可要求 Claude：
- 修改主题（如 `\usetheme{metropolis}`）
- 调整比例（如 `aspectratio=43` 或 `1610/916:9` 等）
- 增加自定义命令和配色方案

---

## 3. 推荐指令示例

### 3.1 从大纲生成 PPT

**示例：**

> 请使用 beamer 帮我生成一份 LaTeX PPT，要求：  
> - 中文支持  
> - 16:9 宽屏  
> - 有封面页和目录页  
> - 根据以下大纲逐页生成：  
>   1. 主题介绍  
>   2. 方法概述（2 页）  
>   3. 实验结果  
>   4. 总结与展望  
> 只需要给出一份**完整可编译**的 `main.tex`。

Claude 预期行为：
- 使用本文件中的基础模板或用户指定模板；
- 每个一级条目对应一到多页 `\begin{frame}...\end{frame}`；
- 自动插入对应的 `\section{}`。

### 3.2 从已有内容转为 PPT

**示例：**

> 下面是我写好的演讲提纲，请帮我：  
> 1. 按逻辑拆分成多页 PPT；  
> 2. 每页使用 `\begin{frame}{标题}`；  
> 3. 要点使用 `itemize` 列表；  
> 4. 最终输出完整的 beamer 文档。  
> （然后粘贴提纲内容）

Claude 预期行为：
- 合理分页，每 3–6 个要点一页；
- 长句改写为简洁 bullet；
- 避免整页大段文字。

---

## 4. 版式与风格约定

1. **页面结构**
   - 每页建议 3–6 条要点；
   - 公式、图表宜单独成页或与少量要点同页；
   - 避免整页长段文本。

2. **常用环境**
   - 列表：`itemize` / `enumerate`
   - 代码：`verbatim` 或 `lstlisting`（如需可加载 `listings` 宏包）
   - 定理类（可选）：
     ```latex
     \begin{theorem}
       ...
     \end{theorem}
     ```

3. **数学公式**
   - 行内公式：`$a^2 + b^2 = c^2$`
   - 重要公式：
     ```latex
     \[
       E = mc^2
     \]
     ```

4. **图片与表格占位符**

图片尚未准备好时，用注释占位：

```latex
\begin{frame}{系统架构}
  % TODO: 在此插入 system_architecture.png
  \begin{itemize}
    \item 模块一：数据采集
    \item 模块二：特征提取
    \item 模块三：模型推理
  \end{itemize}
\end{frame}
```

正式插图示例：

```latex
\begin{frame}{系统架构}
  \begin{figure}
    \centering
    \includegraphics[width=0.8\linewidth]{figs/system_architecture.png}
    \caption{系统整体架构示意图}
  \end{figure}
\end{frame}
```

---

## 5. 编译建议

1. 默认：
   ```bash
   xelatex main.tex
   ```
2. 如有参考文献：
   ```bash
   xelatex main.tex
   bibtex main
   xelatex main.tex
   xelatex main.tex
   ```

---

## 6. Claude 生成代码的规范

1. **必须输出完整文档**
   - 包含 `\documentclass`、宏包、标题信息、`\begin{document}`、`\end{document}`。

2. **禁止省略内容**
   - 不使用“……省略若干页”等占位语；
   - 内容较多时要合理分页并全部给出。

3. **中文支持默认开启**
   - 默认假定编译环境支持 `xelatex`；
   - 若使用 `fontspec` 自定义字体，应使用易替换的字体名（如 `SimSun`、`SimHei` 等）。

4. **文件命名建议**
   - 主文件：`main.tex`
   - 图片目录：`figs/`
   - 若拆分多个 `.tex` 文件，可在输出中说明拆分方式，但这里至少需给出可编译的 `main.tex`。

---

## 7. 完整示例模板

```latex
\documentclass[aspectratio=169,11pt]{beamer}
\usepackage[UTF8]{ctex}

\usetheme{Madrid}
\usecolortheme{dolphin}

\usepackage{amsmath, amssymb, amsfonts}
\usepackage{graphicx}
\usepackage{hyperref}
\usepackage{booktabs}

\title{用 LaTeX 制作 PPT}
\subtitle{Beamer 与工作流示例}
\author{你的名字}
\institute{你的单位}
\date{\today}

\begin{document}

\begin{frame}
  \titlepage
\end{frame}

\begin{frame}{目录}
  \tableofcontents
\end{frame}

\section{为什么使用 LaTeX 做 PPT}
\begin{frame}{优点概览}
  \begin{itemize}
    \item 与论文共享同一套公式与排版宏包
    \item 版本可控，适合放在 Git 仓库中协作
    \item 对数学公式和代码展示非常友好
  \end{itemize}
\end{frame}

\section{基本结构}
\begin{frame}{Beamer 文档结构}
  \begin{itemize}
    \item 使用 \texttt{\textbackslash documentclass\{beamer\}} 声明文类
    \item 在导言区设置主题、字体和常用宏包
    \item 每一页使用 \texttt{frame} 环境定义
  \end{itemize}
\end{frame}

\section{示例}
\begin{frame}{公式示例}
  经典的勾股定理：
  \[
    a^2 + b^2 = c^2
  \]
\end{frame}

\begin{frame}{图片示例}
  \begin{figure}
    \centering
    % \includegraphics[width=0.7\linewidth]{figs/example.png}
    \caption{在此处插入示意图（example.png）}
  \end{figure}
\end{frame}

\section{总结}
\begin{frame}{总结}
  \begin{itemize}
    \item LaTeX + Beamer 适合技术类与学术类汇报
    \item 通过模板和宏命令可以大幅提升制作效率
    \item 与 Git 结合可实现完整的 PPT 版本管理
  \end{itemize}
\end{frame}

\end{document}
```

---
