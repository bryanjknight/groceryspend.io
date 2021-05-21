import React, { useState } from "react";

export interface EditableCellProps<T> {
  id: string;
  value: T;
  valueLabelMaker: (t: T) => string;
  editorFactory: (
    handleChange: (t: T) => void,
    handleOnBlur: () => void,
    defaultValue: T,
  ) => JSX.Element;
  onValueChange: (t: T) => void;
  className: string;
}

// TODO: handle generics better here
// eslint-disable-next-line @typescript-eslint/ban-types
export const EditableCell = <T extends object>(
  props: EditableCellProps<T>
): JSX.Element => {
  const [value, setValue] = useState(null as unknown as T);
  const [editing, setEditing] = useState(false);

  if (editing) {
    const handleChange = (t: T) => {
      setValue(t);
      setEditing(false);
      props.onValueChange(t);
    };

    const handleOnBlur = () => {
      setEditing(false);
    };

    const editor = props.editorFactory(handleChange, handleOnBlur, value || props.value);

    return <td className={props.className}>{editor}</td>;
  } else {
    return (
      <td onClick={() => setEditing(true)} className={props.className} style={{cursor: "pointer"}}>
        {props.valueLabelMaker(value || props.value)}
      </td>
    );
  }
};
