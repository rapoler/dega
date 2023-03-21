import React, { useState } from 'react';
import { Layout, Card, notification, BackTop, ConfigProvider } from 'antd';
import { withRouter } from 'react-router-dom';
import { Helmet } from 'react-helmet';
import Sidebar from '../components/GlobalNav/Sidebar';
import Header from '../components/GlobalNav/Header';
import { useDispatch, useSelector } from 'react-redux';
import { getSpaces } from '../actions/spaces';
import './basic.css';
import { getSuperOrganisation } from '../actions/admin';
import Pageheader from '../components/PageHeader';
import routes from '../config/routesConfig';
import _ from 'lodash';

function BasicLayout(props) {
  const { location } = props;
  const { Content } = Layout;
  const { children } = props;
  const [enteredRoute, setRoute] = React.useState({ menuKey: '/' });
  React.useEffect(() => {
    const pathSnippets = location.pathname.split('/').filter((i) => i);
    if (pathSnippets.length === 0) {
      setRoute({ menuKey: '/' });
      return;
    }
    var index;
    for (index = 0; index < pathSnippets.length; index++) {
      const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
      const nextTempRoute =
        pathSnippets.length - index > 1
          ? _.find(routes, { path: `/${pathSnippets.slice(0, index + 2).join('/')}` })
          : null;
      const tempRoute = _.find(routes, { path: url });
      if (nextTempRoute) {
        continue;
      }
      if (tempRoute) {
        setRoute(tempRoute);
        break;
      }
    }
  }, [location]);
  const dispatch = useDispatch();

  const { permission, orgs, loading, selected, applications, services } = useSelector((state) => {
    const { selected, orgs, loading } = state.spaces;

    if (selected > 0) {
      const space = state.spaces.details[selected];

      const applications = orgs.find((org) => org.spaces.includes(space.id))?.applications || [];

      return {
        applications: applications,
        permission: space.permissions || [],
        orgs: orgs,
        loading: loading,
        selected: selected,
        services: space.services,
      };
    }
    return {
      orgs: orgs,
      loading: loading,
      permission: [],
      selected: selected,
      applications: [],
      services: ['core'],
    };
  });

  const { type, message, description, time, redirect } = useSelector((state) => {
    return { ...state.notifications, redirect: state.redirect };
  });

  const superOrg = useSelector(({ admin }) => {
    return admin.organisation;
  });

  React.useEffect(() => {
    dispatch(getSpaces());
    // .then((org) => {
    // if (org && org.length > 0) dispatch(getSuperOrganisation(org[0].id));
    // }
    // );
  }, [dispatch, selected]);

  React.useEffect(() => {
    if (type && message && description && selected !== 0) {
      notification[type]({
        message: message,
        description: description,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [time, description]);

  React.useEffect(() => {
    if (redirect?.code === 307) {
      window.location.href = window.REACT_APP_KAVACH_PUBLIC_URL;
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [redirect]);

  // React.useEffect(() => {
  //   if (orgs.length > 0 && selected === 0) history.push('/spaces/create');
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, [orgs, location.pathname]);

  const hideSidebar =
    (location.pathname.includes('posts') ||
      location.pathname.includes('fact-checks') ||
      location.pathname.includes('pages')) &&
    (location.pathname.includes('edit') || location.pathname.includes('create'));
  return (
    <ConfigProvider
    // theme={
    // {
    // we can customize the theme here , commenting it for now.
    // token: {
    //   colorPrimary: '#4F46E5',
    //   colorLink: '#4F46E5',
    // },
    // components: {
    //   Menu: {
    //     colorItemBgSelected: '#D1D5DB',
    //   },
    // },
    // }
    // }
    >
      <Layout hasSider={true}>
        <Helmet titleTemplate={'%s | Dega Studio'} title={'Dega Studio'} />
        {!hideSidebar && (
          <Sidebar
            permission={permission}
            menuKey={enteredRoute?.menuKey}
            orgs={orgs}
            loading={loading}
            superOrg={superOrg}
            applications={applications}
            services={services}
          />
        )}
        <Layout style={{ background: '#fff' }}>
          {/* <Header applications={applications} hideSidebar={hideSidebar} /> */}
          <Content className="layout-content">
            <Pageheader location={location} />
            <Card key={selected.toString()} className="wrap-children-content">
              {children}
            </Card>
          </Content>
          <BackTop style={{ right: 50 }} />
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}

export default withRouter(BasicLayout);
