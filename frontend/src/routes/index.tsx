import React from "react";
import PathConstants from "./PathConstants";
import ProtectedRoute from "./../components/common/Protected";

const Home = React.lazy(() => import("../components/pages/home/HomePage"));
const Login = React.lazy(() => import("../components/pages/login/LoginPage"));

const routes = [
  {
    path: PathConstants.HOME,
    element: (
      <ProtectedRoute>
        <Home />
      </ProtectedRoute>
    ),
  },
  {path: PathConstants.LOGIN, element: <Login />},
];

export default routes;
