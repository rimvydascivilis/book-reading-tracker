import {createBrowserRouter, RouterProvider} from 'react-router-dom';
import routes from './routes';
import AppLayout from './components/Applayout/AppLayout';
import React from 'react';

const NotFoundPage = React.lazy(() => import('./components/pages/404Page'));

const App: React.FC = () => {
  const router = createBrowserRouter([
    {
      element: <AppLayout />,
      errorElement: (
        <AppLayout>
          <NotFoundPage />
        </AppLayout>
      ),
      children: routes,
    },
  ]);

  return <RouterProvider router={router} />;
};

export default App;
