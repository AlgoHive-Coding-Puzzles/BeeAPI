import { useState, useEffect } from "react";
import { getAllowedElements } from "../utils/htmlConverter";

interface HtmlPreviewProps {
  html: string;
}

const HtmlPreview = ({ html }: HtmlPreviewProps) => {
  const [sanitizedHtml, setSanitizedHtml] = useState("");
  const allowedElements = getAllowedElements();

  useEffect(() => {
    // Simple sanitization - in a real app you would use a proper sanitizer library
    const temp = document.createElement("div");
    temp.innerHTML = html;

    // Remove disallowed tags
    const allElements = temp.getElementsByTagName("*");
    for (let i = allElements.length - 1; i >= 0; i--) {
      const el = allElements[i];
      if (!allowedElements.includes(el.tagName.toLowerCase())) {
        el.parentNode?.removeChild(el);
      }
    }

    // Remove any script tags for security
    const scripts = temp.getElementsByTagName("script");
    while (scripts[0]) {
      scripts[0].parentNode?.removeChild(scripts[0]);
    }

    // Remove all style and class attributes
    const allElementsWithAttrs = temp.getElementsByTagName("*");
    for (let i = 0; i < allElementsWithAttrs.length; i++) {
      const el = allElementsWithAttrs[i];
      el.removeAttribute("style");
      // Don't remove class attributes as we might need them for styling
    }

    // Process code elements to fix whitespace issues
    const codeElements = temp.querySelectorAll("code");
    codeElements.forEach((code) => {
      // Get the text content and split into lines
      const content = code.textContent || "";
      if (content.trim()) {
        const lines = content.split("\n");

        // Find the minimum indentation across all non-empty lines
        let minIndent = Infinity;
        lines.forEach((line) => {
          const trimmedLine = line.trimStart();
          if (trimmedLine.length > 0) {
            const indent = line.length - trimmedLine.length;
            minIndent = Math.min(minIndent, indent);
          }
        });

        // Apply the indentation correction to all lines
        const correctedLines = lines.map((line) => {
          if (line.trim().length === 0) return line;
          return line.slice(minIndent === Infinity ? 0 : minIndent);
        });

        code.textContent = correctedLines.join("\n");
        code.classList.add("preview-code-block");
      }
    });

    setSanitizedHtml(temp.innerHTML);
  }, [allowedElements, html]);

  return (
    <div
      className="html-preview w-full overflow-auto border rounded-md p-4 bg-white"
      style={{ minHeight: "500px" }}
    >
      <style>{`
        .preview-content code, .preview-code-block {
          white-space: pre-wrap;
          display: block;
          background-color: #f5f5f5;
          padding: 0.75rem;
          border-radius: 4px;
          overflow-x: auto;
          font-family: 'Courier New', Courier, monospace;
          color: #333;
        }
        .preview-content h2, .preview-content h3, .preview-content h4, 
        .preview-content h5, .preview-content h6 {
          margin-top: 16px;
          margin-bottom: 8px;
          font-weight: bold;
          color: #333;
        }
        .preview-content p {
          margin: 8px 0;
          color: #333;
        }
      `}</style>
      <div
        dangerouslySetInnerHTML={{ __html: sanitizedHtml }}
        className="preview-content"
      />
    </div>
  );
};

export default HtmlPreview;
