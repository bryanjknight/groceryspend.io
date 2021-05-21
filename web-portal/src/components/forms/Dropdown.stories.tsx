import React from 'react';
import { Story, Meta } from '@storybook/react';

import { Dropdown, DropdownProps } from './Dropdown';
import { Category } from '../../models';

export default {
  title: 'Dropdown',
  component: Dropdown,
  argTypes: {
    backgroundColor: { control: 'color' },
  },
} as Meta;

const Template: Story<DropdownProps<Category>> = (args) => <Dropdown {...args} />;

export const NoSelectedValue = Template.bind({});
NoSelectedValue.args = {
  id: "test",
  mapOptionsToSelectItems: (c: Category) => ({label: c.Name, value: c.ID.toString()}),
  onSelect: (c: Category) => console.log(`Selected: ${c.Name}`),
  options: [
    {
      ID: 1,
      Name: "A",
    },
    {
      ID: 2,
      Name: "B",
    },
    {
      ID: 3,
      Name: "C",
    },
  ]
};

export const DefaultValue = Template.bind({});
DefaultValue.args = {
  id: "test-default-value",
  mapOptionsToSelectItems: (c: Category) => ({label: c.Name, value: c.ID.toString()}),
  onSelect: (c: Category) => console.log(`Selected: ${c.Name}`),
  options: [
    {
      ID: 1,
      Name: "A",
    },
    {
      ID: 2,
      Name: "B",
    },
    {
      ID: 3,
      Name: "C",
    },
  ],
  defaultValue: {
    ID: 1,
    Name: "A",
  },
};