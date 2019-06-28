/* ************************************************************************** */
/*                                                                            */
/*  UserList.js                                                               */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Jun 28 15:19:57 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

import React, { useState, useCallback, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCoins } from '@fortawesome/free-solid-svg-icons';
import { List, Avatar, Row, Col, message } from 'antd';
import { JsonRequest } from '../common/DataHandler';
import { Config } from '../common/config';

const UserDetail = ({ coins }) => {
  return (
    <Row type="flex" justify="start" align="start" gutter={10}>
      <Col>
        {coins} <FontAwesomeIcon icon={faCoins} />
      </Col>
      <Col>| (n) PR reviewed</Col>
    </Row>
  );
};

const UserList = () => {
  const [userData, setUserData] = useState(null);

  const getDataSource = useCallback(() => {
    return userData ? userData.map(user => ({
        ...user,
        content: <UserDetail coins={user.coins}/>
    })) : [];
  }, [userData]);

  const getUserData = useCallback(() => {
    const onsuccess = response => {
      const data = (response && JSON.parse(response)) || [];
      message.success('user data updated successfully');
      setUserData(data);
    };

    const onerror = () => {
      const user = this.props.store.user || {};

      this.setState({
        user,
        prevuser: user
      });
      message.error('user data failed to update');
    };

    JsonRequest({
      url: `http://${Config.apiDomainName}/users`,
      method: 'get',
      onsuccess,
      onerror
    });
  }, [userData]);

  useEffect(() => {
      if (userData === null) {
        getUserData();
      }
  });

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
      dataSource={getDataSource()}
      renderItem={item => (
        <List.Item>
          <List.Item.Meta
            avatar={
              <Avatar src={item.avatar_url} size="large"/>
            }
            title={<a href={item.url}>{item.login}</a>}
          />
          {item.content}
        </List.Item>
      )}
    />
  );
};

export default UserList;
