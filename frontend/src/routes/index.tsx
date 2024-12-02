import React from 'react';
import PathConstants from './PathConstants';
import ProtectedRoute from './../components/common/Protected';
import {Navigate} from 'react-router-dom';

const Login = React.lazy(() => import('../components/pages/login/LoginPage'));
const Library = React.lazy(
  () => import('../components/pages/library/LibraryPage'),
);
const Goal = React.lazy(() => import('../components/pages/goal/GoalPage'));
const MyReads = React.lazy(
  () => import('../components/pages/my-reads/MyReadsPage'),
);
const Lists = React.lazy(() => import('../components/pages/lists/ListsPage'));

const protectedRoutes = [
  {
    path: PathConstants.LIBRARY,
    component: Library,
  },
  {
    path: PathConstants.GOAL,
    component: Goal,
  },
  {
    path: PathConstants.MY_READS,
    component: MyReads,
  },
  {
    path: PathConstants.LISTS,
    component: Lists,
  },
];

const routes = [
  // unauthenticated routes
  {path: '/', element: <Navigate to={PathConstants.HOME} />},
  {path: PathConstants.LOGIN, element: <Login />},

  ...protectedRoutes.map(route => ({
    path: route.path,
    element: (
      <ProtectedRoute>
        <route.component />
      </ProtectedRoute>
    ),
  })),
];

export default routes;
