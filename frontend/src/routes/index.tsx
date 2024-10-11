import React from "react";
import PathConstants from "./PathConstants";

const Home = React.lazy(() => import("../components/pages/home/Home"));

const routes = [
  { path: PathConstants.HOME, element: <Home /> },
];

export default routes;