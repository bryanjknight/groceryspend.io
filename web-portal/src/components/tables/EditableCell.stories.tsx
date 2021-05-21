import React from "react";
import { Story, Meta } from "@storybook/react";

import { EditableCell, EditableCellProps } from "./EditableCell";
import { Category } from "../../models";
import { Dropdown, DropdownProps } from "../forms/Dropdown";

export default {
  title: "EditableCell",
  component: EditableCell,
  argTypes: {
    backgroundColor: { control: "color" },
  },
} as Meta;

const Template: Story<EditableCellProps<Category>> = (args) => (
  <table className="table">
    <thead>
      <tr className="d-flex">
        <th className="col-1">Col A</th>
        <th className="col-1">Col B</th>
      </tr>
    </thead>
    <tbody>
      <tr className="d-flex">
        <td className="col-1">Not editable cell</td>
        <EditableCell {...args} className="col-1" />
      </tr>
    </tbody>
  </table>
);

export const Primary = Template.bind({});

Primary.args = {
  id: "test",
  editorFactory: (handleChange, handleOnBlur) => {
    const dropdownProps: DropdownProps<Category> = {
      id: "test-dropdown",
      mapOptionsToSelectItems: (c: Category) => ({
        label: c.Name,
        value: c.ID.toString(),
      }),
      onSelect: (c: Category) => handleChange(c),
      onBlur: () => handleOnBlur(),
      options: [
        { ID: 1, Name: "A" },
        { ID: 2, Name: "B" },
        { ID: 3, Name: "C" },
      ],
    };
    return <Dropdown {...dropdownProps} />
  },
  valueLabelMaker: (c: Category) => c.Name,
  value: {
    ID: 1,
    Name: "A",
  },
};
