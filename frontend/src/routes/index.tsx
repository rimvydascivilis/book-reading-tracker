import React from "react";
import PathConstants from "./PathConstants";
import ProtectedRoute from "./../components/common/Protected";

const Login = React.lazy(() => import("../components/pages/login/LoginPage"));
const Library = React.lazy(
  () => import("../components/pages/library/LibraryPage"),
);

const routes = [
  {
    path: PathConstants.LIBRARY,
    element: (
      <ProtectedRoute>
        <Library />
      </ProtectedRoute>
    ),
  },
  {path: PathConstants.LOGIN, element: <Login />},
];

export default routes;
