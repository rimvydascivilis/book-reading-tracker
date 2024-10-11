import { createBrowserRouter, RouterProvider } from "react-router-dom";
import routes from "./routes";
import ErrorPage from "./components/pages/ErrorPage";
import AppLayout from "./components/layouts/AppLayout";

export default function App() {
  const router = createBrowserRouter([
    {
      element: <AppLayout />,
      errorElement: <ErrorPage />,
      children: routes,
    },
  ]);

  return <RouterProvider router={router} />;
}