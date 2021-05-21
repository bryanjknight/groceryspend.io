import React from "react";

export interface DropdownProps<T> {
  id: string;
  onSelect: (t: T) => void;
  mapOptionsToSelectItems: (t: T) => OptionType;
  options: T[];
  defaultValue?: T;
}

export interface OptionType {
  value: string;
  label: string;
}

export type DropdownJSXElement = <T>(props: DropdownProps<T>) => JSX.Element;

export const Dropdown: DropdownJSXElement = (props): JSX.Element => {

  const valueToOption = props.options.reduce((acc, option) => {
    const selectOption = props.mapOptionsToSelectItems(option)
    return {
      ...acc,
      [selectOption.value]: option,
    }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  }, {} as Record<string, any>); // TODO: is there a better way to handle the generics


  
  const options = props.options
    .map(props.mapOptionsToSelectItems)
    .map((item) => (
      <option value={item.value}>
        {item.label}
      </option>
    ));
  
  if (!props.defaultValue) {
    options.splice(0, 0, <option value="">--Select One---</option>)
  }
  
  let extraArgs = {};
  if (props.defaultValue) {
    extraArgs = {
      ...extraArgs,
      defaultValue: props.defaultValue,
    }
  }

  const handleOnSelect = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const selectValue = event.target.value

    if (selectValue === "") return;

    const item = valueToOption[selectValue];

    props.onSelect(item);
  }

  return (
    <div>
      <select id={props.id} onSelect={handleOnSelect} {...extraArgs}>{options}</select>
    </div>
  );
};
