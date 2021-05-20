import React from "react";

export interface DropdownProps<T> {
  id: string;
  onSelect: (t: T) => void;
  mapOptionsToSelectItems: (t: T) => OptionType;
  options: T[];
  selectedValue?: T;
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

  const selectedValueID = props.selectedValue
    ? props.mapOptionsToSelectItems(props.selectedValue)
    : null;
  const defaultValueID = props.defaultValue
    ? props.mapOptionsToSelectItems(props.defaultValue)
    : null;

  const isSelected = (item: OptionType) => {
    if ((selectedValueID && selectedValueID.value === item.value) ||
    (defaultValueID && defaultValueID.value === item.value)) {
      return true;
    }
    return false;
  }

  const options = props.options
    .map(props.mapOptionsToSelectItems)
    .map((item) => (
      <option value={item.value} selected={isSelected(item)}>
        {item.label}
      </option>
    ));

  const handleOnSelect = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const selectValue = event.target.value

    const item = valueToOption[selectValue];

    props.onSelect(item);
  }

  return (
    <div>
      <select id={props.id} onSelect={handleOnSelect}>{options}</select>
    </div>
  );
};
