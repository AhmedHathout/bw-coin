/* ************************************************************************** */
/*                                                                            */
/*  App.js                                                                    */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created: Fri Jun 28 10:17:56 2019                        by elhmn        */
/*   Updated: Fri Jun 28 10:25:26 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

import React from 'react';
import { BrowserRouter as Router, Route } from "react-router-dom";
import Home from './components/Home';
import "antd/dist/antd.css";

function App() {
  return (
      <Router>
        <Route path="/" exact component={Home} />
      </Router>
  );
}

export default App;
