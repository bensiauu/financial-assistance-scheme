import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ApplicationList from './pages/ApplicationList'; // Your ApplicationList component
import SchemeList from './pages/SchemeList'; // Your SchemeList component
import ApplicantList from './pages/ApplicantList'; // Your ApplicantList component
import ApplicationDetails from './pages/ApplicationDetails'; // Your ApplicationDetails component
import MainLayout from './MainLayout'; // MainLayout that includes the sidebar

function App() {
  return (
    <Router>
      <Routes>
        <Route element={<MainLayout />}>
          <Route path="/applications" element={<ApplicationList />} />
          <Route path="/schemes" element={<SchemeList />} />
          <Route path="/applicants" element={<ApplicantList />} />
          <Route path="/applications/:id" element={<ApplicationDetails />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
