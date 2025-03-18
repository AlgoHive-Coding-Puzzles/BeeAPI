import { Button } from "primereact/button";
import { useEffect, useState } from "react";
import HtmlEditor from "./HtmlEditor";
import RichTextContentEditable from "./RichTextContentEditable";
import {
  htmlToRichText,
  richTextToHtml,
  applyEditorClasses,
  removeEditorClasses,
} from "../utils/htmlConverter";

interface HtmlRichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
}

const HtmlRichTextEditor = ({ value, onChange }: HtmlRichTextEditorProps) => {
  const [editorMode, setEditorMode] = useState<"rich" | "html">("rich");
  const [internalValue, setInternalValue] = useState(value);

  useEffect(() => {
    // Initialize with the provided value
    setInternalValue(value);
  }, [value]);

  const handleChange = (newValue: string) => {
    setInternalValue(newValue);
    onChange(editorMode === "rich" ? richTextToHtml(newValue) : newValue);
  };

  const toggleEditorMode = () => {
    if (editorMode === "rich") {
      // Switch to HTML mode - remove editor classes before showing HTML
      const cleanHtml = removeEditorClasses(internalValue);
      setInternalValue(cleanHtml);
      setEditorMode("html");
    } else {
      // Switch to rich text mode - convert HTML to rich text format and add classes
      const richText = htmlToRichText(internalValue);
      const formattedRichText = applyEditorClasses(richText);
      setInternalValue(formattedRichText);
      setEditorMode("rich");
    }
  };

  return (
    <div className="html-rich-text-editor w-full">
      <div className="editor-mode-toggle mb-2 flex justify-content-end">
        <Button
          label={editorMode === "rich" ? "Show HTML" : "Rich Text Mode"}
          icon={editorMode === "rich" ? "pi pi-code" : "pi pi-eye"}
          onClick={toggleEditorMode}
          className="p-button-sm"
        />
      </div>

      {editorMode === "rich" ? (
        <RichTextContentEditable
          value={internalValue}
          onChange={handleChange}
        />
      ) : (
        <HtmlEditor value={internalValue} onChange={handleChange} />
      )}
    </div>
  );
};

export default HtmlRichTextEditor;
