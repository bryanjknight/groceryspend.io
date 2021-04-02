import React, { MouseEventHandler } from "react";

export interface ButtonProps {
  text: string;
  dataTestId?: string;
  onClick: MouseEventHandler<HTMLButtonElement>;
}

export const Button = (props: ButtonProps): JSX.Element => {
  return (
    <button onClick={props.onClick} data-testid={props.dataTestId || ""}>
      {props.text}
    </button>
  );
};
