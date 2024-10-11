import {
  useNavigate,
  isRouteErrorResponse,
  useRouteError,
} from 'react-router-dom';
import { Button } from "antd";

export default function ErrorPage() {
  const navigate = useNavigate();
  const error = useRouteError() as Error;

  if (!isRouteErrorResponse(error)) {
    return null;
  }

  return (
    <div>
      <h1>Something went wrong 😢</h1>
      <p>{error.data}</p>
      <Button onClick={() => navigate(-1)}> Go back</Button>
    </div>
  );
};