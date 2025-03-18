/**
 * Converts raw HTML to a format suitable for the rich text editor
 */
export const htmlToRichText = (html: string): string => {
  // If the HTML is empty, provide a minimal structure
  if (!html || html.trim() === "") {
    return "<article><p></p></article>";
  }

  // Create a DOM parser to handle the HTML
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, "text/html");

  // Check if there's an article tag
  const articleElement = doc.querySelector("article");

  if (articleElement) {
    // Return the full article with its tags for contenteditable
    return articleElement.outerHTML;
  } else {
    // Wrap the content in an article tag
    return `<article>${html}</article>`;
  }
};

/**
 * Converts rich text editor content to proper HTML with article tags
 */
export const richTextToHtml = (richText: string): string => {
  // Check if content already has article tags
  if (
    richText.trim().startsWith("<article") &&
    richText.trim().endsWith("</article>")
  ) {
    return richText;
  }

  // If content doesn't have article tags, add them
  return `<article>${richText}</article>`;
};

/**
 * Ensures the HTML has an article tag as the root
 */
export const ensureArticleTag = (html: string): string => {
  if (
    html.trim().startsWith("<article") &&
    html.trim().endsWith("</article>")
  ) {
    return html;
  }

  // If no article tag, wrap the content
  return `<article>${html}</article>`;
};

/**
 * Gets the allowed HTML elements for the rich text editor
 */
export const getAllowedElements = (): string[] => {
  return [
    "article",
    "h2",
    "h3",
    "h4",
    "h5",
    "h6",
    "p",
    "span",
    "div",
    "code",
    "pre",
    "br",
  ];
};

/**
 * Get all the block for the rich text editor
 */
export const getBlockElements = (): { tag: string; label: string }[] => {
  return [
    {
      tag: "h2",
      label: "Heading 2",
    },
    {
      tag: "h3",
      label: "Heading 3",
    },
    {
      tag: "h4",
      label: "Heading 4",
    },
    {
      tag: "h5",
      label: "Heading 5",
    },
    {
      tag: "h6",
      label: "Heading 6",
    },
    {
      tag: "p",
      label: "Paragraph",
    },
    {
      tag: "code",
      label: "Code Block",
    },
  ];
};

/**
 * Creates an empty template with the allowed elements
 */
export const createEditorTemplate = (): string => {
  return `<article>
    <h2>HTML Editor</h2>
    <p>Start typing your HTML here...</p>
    <code>
    1
    2
    3
    </code>
    <h3>Preview</h3>
    <p>This is a preview of your HTML.</p>
    </article>`;
};

/**
 * Apply specific CSS classes to HTML elements for the rich text editor
 */
export const applyEditorClasses = (htmlContent: string): string => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(htmlContent, "text/html");

  // Add classes to specific elements
  doc.querySelectorAll("code").forEach((el) => {
    el.classList.add("editor-code-block");
  });

  doc.querySelectorAll("h2, h3, h4, h5, h6").forEach((el) => {
    el.classList.add("editor-heading");
  });

  doc.querySelectorAll("p").forEach((el) => {
    el.classList.add("editor-paragraph");
  });

  const body = doc.querySelector("body");
  return body ? body.innerHTML : htmlContent;
};

/**
 * Removes editor-specific classes from HTML content
 */
export const removeEditorClasses = (htmlContent: string): string => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(htmlContent, "text/html");

  // Delete all classes
  doc.querySelectorAll("*").forEach((el) => {
    el.removeAttribute("class");
  });

  const body = doc.querySelector("body");
  return body ? body.innerHTML : htmlContent;
};
