import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Login from "./pages/Login";
import ApplicationList from "./pages/ApplicationList";
import ApplicationDetails from "./pages/ApplicationDetails";
import FormSubmission from "./pages/FormSubmission";
import axios from "axios";



function App() {
  axios.defaults.withCredentials = true;
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/applications" element={<ApplicationList />} />
        <Route path="/applications/:id" element={<ApplicationDetails />} />
        <Route path="/apply" element={<FormSubmission />} />
      </Routes>
    </Router>
  );
}


export default App;
