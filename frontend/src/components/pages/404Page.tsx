import {
  useNavigate,
  isRouteErrorResponse,
  useRouteError,
} from "react-router-dom";
import {Button} from "antd";
import React from "react";

const NotFoundPage: React.FC = () => {
  const navigate = useNavigate();
  const error = useRouteError() as Error;

  if (!isRouteErrorResponse(error)) {
    return null;
  }

  return (
    <div style={{textAlign: "center"}}>
      <h1>Something went wrong 😢</h1>
      <p>{error.data}</p>
      <Button onClick={() => navigate(-1)}> Go back</Button>
    </div>
  );
};

export default NotFoundPage;
