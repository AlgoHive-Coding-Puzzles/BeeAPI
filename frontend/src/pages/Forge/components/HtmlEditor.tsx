import { useState, useEffect } from "react";
import { InputTextarea } from "primereact/inputtextarea";

interface HtmlEditorProps {
  value: string;
  onChange: (value: string) => void;
}

const HtmlEditor = ({ value, onChange }: HtmlEditorProps) => {
  const [internalValue, setInternalValue] = useState(value);

  useEffect(() => {
    setInternalValue(value);
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    setInternalValue(newValue);
    onChange(newValue);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Tab") {
      e.preventDefault();

      // Get cursor position
      const target = e.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;

      // Insert tab at cursor position
      const newValue =
        internalValue.substring(0, start) + "\t" + internalValue.substring(end);

      // Update state
      setInternalValue(newValue);
      onChange(newValue);

      // Reset cursor position after tab insertion
      setTimeout(() => {
        target.selectionStart = target.selectionEnd = start + 1;
      }, 0);
    }
  };

  return (
    <div className="html-editor w-full h-full">
      <InputTextarea
        value={internalValue}
        onChange={handleChange}
        onKeyDown={handleKeyDown}
        className="w-full"
        style={{
          minHeight: "500px",
          fontFamily: "monospace",
          fontSize: "14px",
          lineHeight: 1.5,
        }}
        autoResize
      />
    </div>
  );
};

export default HtmlEditor;
