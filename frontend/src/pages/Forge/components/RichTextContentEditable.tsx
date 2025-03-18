import { useEffect, useRef } from "react";
import {
  applyEditorClasses,
  createEditorTemplate,
  getBlockElements,
} from "../utils/htmlConverter";

interface RichTextContentEditableProps {
  value: string;
  onChange: (value: string) => void;
}

const RichTextContentEditable = ({
  value,
  onChange,
}: RichTextContentEditableProps) => {
  const editorRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (editorRef.current) {
      // Only update the content if it's different to avoid losing cursor position
      if (editorRef.current.innerHTML !== value) {
        // Apply editor classes when initializing content
        const formattedValue = applyEditorClasses(value);
        editorRef.current.innerHTML = formattedValue;
      }
    }
  }, [value]);

  const handleInput = () => {
    if (editorRef.current) {
      onChange(editorRef.current.innerHTML);
    }
  };

  // Insert placeholder if the editor is empty
  const handleBlur = () => {
    if (editorRef.current && !editorRef.current.innerHTML.trim()) {
      editorRef.current.innerHTML = createEditorTemplate();
      onChange(editorRef.current.innerHTML);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Tab") {
      e.preventDefault();
      document.execCommand("insertHTML", false, "&nbsp;&nbsp;&nbsp;&nbsp;");
    }
  };

  // Handle inserting specific elements
  const insertElement = (tagName: string) => {
    const selection = window.getSelection();
    if (!selection || !selection.rangeCount) return;

    console.log(selection);

    // If the selection is not {anchorNode: text
    if (selection.anchorNode?.nodeName !== "#text") {
      selection.collapseToEnd();
    }

    const range = selection.getRangeAt(0);
    const element = document.createElement(tagName);

    if (tagName === "code") {
      element.textContent = "// Your code here";
      element.className = "editor-code-block";
    } else {
      element.textContent =
        tagName === "p" ? "New paragraph" : `New ${tagName} heading`;
      if (tagName.startsWith("h")) {
        element.className = "editor-heading";
      } else if (tagName === "p") {
        element.className = "editor-paragraph";
      }
    }

    range.deleteContents();
    range.insertNode(element);

    // Position cursor at the end of the new element
    range.setStartAfter(element);
    range.setEndAfter(element);
    selection.removeAllRanges();
    selection.addRange(range);

    // Notify change
    handleInput();
  };

  return (
    <div className="rich-text-editor-container">
      <div className="editor-toolbar">
        {getBlockElements().map((el) => (
          <button
            key={el.tag}
            type="button"
            className="toolbar-button"
            onClick={() => insertElement(el.tag)}
          >
            {el.label}
          </button>
        ))}
      </div>
      <div
        ref={editorRef}
        className="rich-text-editor"
        contentEditable={true}
        onInput={handleInput}
        onBlur={handleBlur}
        onKeyDown={handleKeyDown}
        style={{
          minHeight: "500px",
          border: "1px solid #ced4da",
          borderRadius: "4px",
          padding: "10px",
          outline: "none",
          overflowY: "auto",
          fontFamily: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif",
          lineHeight: "1.5",
        }}
      />
      <style>{`
        .rich-text-editor {
          position: relative;
        }
        .editor-code-block {
          display: block;
          background-color: #1e1e1e;
          padding: 12px;
          border-radius: 4px;
          font-family: "Courier New", Courier, monospace;
          white-space: pre-wrap;
          margin: 10px 0;
        }
        .editor-heading {
          margin: 16px 0 8px;
          font-weight: bold;
        }
        .editor-paragraph {
          margin: 8px 0;
        }
        .editor-toolbar {
          margin-bottom: 8px;
          display: flex;
          gap: 4px;
        }
        .toolbar-button {
          padding: 4px 8px;
          background-color: #ff690b;
          border: 1px solid #414448;
          border-radius: 4px;
          cursor: pointer;
        }
        .toolbar-button:hover {
          background-color: #b24907;
        }
      `}</style>
    </div>
  );
};

export default RichTextContentEditable;
