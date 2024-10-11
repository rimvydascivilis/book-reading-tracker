import { render, screen } from "@testing-library/react";
import { act } from 'react';
import App from "./App";


describe("App Component", () => {
  afterEach(() => {
    jest.clearAllMocks();
    jest.clearAllTimers();
  });

  it("should render AppLayout", async () => {
    await act(async () => {
      render(<App />);
    });

    expect(screen.getByText(/Book reading tracker/i)).toBeInTheDocument(); // Checking footer text
  });
});
