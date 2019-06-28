/* ************************************************************************** */
/*                                                                            */
/*  UserList.js                                                               */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Jun 28 13:06:33 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCoins } from '@fortawesome/free-solid-svg-icons';
import { List, Avatar, Row, Col } from 'antd';

const UserDetail = () => {
  return (
    <Row type="flex" justify="start" align="start" gutter={10}>
      <Col>
        190 <FontAwesomeIcon icon={faCoins} />
      </Col>
      <Col>| 7 PR reviewed</Col>
    </Row>
  );
};

const data = [
  {
    title: 'Element',
    content: <UserDetail />
  },
  {
    title: 'Boris Mbarga',
    content: <UserDetail />
  }
];

const UserList = () => {
  return (
    <List
      itemLayout="vertical"
      size="large"
      pagination={{
        onChange: page => {
          console.log(page);
        },
        pageSize: 10
      }}
      bordered
      dataSource={data}
      renderItem={item => (
        <List.Item>
          <List.Item.Meta
            avatar={
              <Avatar src="https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png" />
            }
            title={<a href="userurl">{item.title}</a>}
          />
          {item.content}
        </List.Item>
      )}
    />
  );
};

export default UserList;
