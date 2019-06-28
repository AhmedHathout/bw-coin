/* ************************************************************************** */
/*                                                                            */
/*  Home.js                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created: Fri Jun 28 10:20:40 2019                        by elhmn        */
/*   Updated: Fri Jun 28 13:06:53 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

import React from 'react';
import { Row, Col, PageHeader } from 'antd';
import UserList from './UserList';

const { Fragment } = React;

const Home = () => {
  return (
    <Fragment>
      <Row type="flex" justify="space-around" align="middle">
      <PageHeader
        title="Users"
      />
      </Row>
      <Row type="flex" justify="space-around" align="middle">
        <Col span={12}>
            <UserList/>
        </Col>
      </Row>
    </Fragment>
  );
};

export default Home;
