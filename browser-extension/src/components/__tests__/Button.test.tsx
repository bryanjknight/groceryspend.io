import React from "react";
import { render, unmountComponentAtNode } from "react-dom";
import { act } from "react-dom/test-utils";
import { Button } from "../Button";

describe("Button Tests", () => {
  const handler = jest.fn();
  let container: HTMLDivElement | null = null;

  beforeEach(() => {
    container = document.createElement("div");
    document.body.appendChild(container);
    handler.mockClear();
  });

  afterEach(() => {
    if (container) {
      unmountComponentAtNode(container);
      container.remove();
      container = null;
    }
  });

  it("Renders text when supplied text", () => {
    const text = "Hello World";
    act(() => {
      render(<Button onClick={handler} text={text} />, container);
    });
    expect(container?.textContent).toBe(text);
  });

  it("Triggers callback when clicked", () => {
    const text = "Hello World";
    act(() => {
      render(
        <Button onClick={handler} text={text} dataTestId="blah" />,
        container
      );
    });
    const button = document.querySelector("[data-testid=blah]");
    if (!button) {
      throw new Error("Button not found")
    }
    expect(button.innerHTML).toBe(text);

    act(() => {
      button.dispatchEvent(new MouseEvent("click", { bubbles: true }));
    });
    expect(handler).toHaveBeenCalledTimes(1);
  });
});
