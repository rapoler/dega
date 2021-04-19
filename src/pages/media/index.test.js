import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { useDispatch, useSelector, Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import thunk from 'redux-thunk';
import { mount } from 'enzyme';

import '../../matchMedia.mock';
import Media from './index';

const middlewares = [thunk];
const mockStore = configureMockStore(middlewares);

jest.mock('react-redux', () => ({
  ...jest.requireActual('react-redux'),
  useDispatch: jest.fn(),
  useSelector: jest.fn(),
}));

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: jest.fn(),
}));

jest.mock('../../actions/media', () => ({
  getMedia: jest.fn(),
  addMedium: jest.fn(),
}));

describe('Media List component', () => {
  let store;
  let mockedDispatch;

  beforeEach(() => {
    store = mockStore({});
    store.dispatch = jest.fn(() => ({}));
    mockedDispatch = jest.fn();
    useDispatch.mockReturnValue(mockedDispatch);
  });
  it('should render the component', () => {
    useSelector.mockImplementation(() => ({}));
    const tree = mount(
      <Provider store={store}>
        <Router>
          <Media permission={{ actions: ['create'] }} />
        </Router>
      </Provider>,
    )
    expect(tree).toMatchSnapshot();
  });
  it('should render the component with data', () => {
    useSelector.mockImplementation(() => ({
      media: [
        {
          id: 1,
          name: 'name',
          url: 'some-url',
          file_size: 'file_size',
          caption: 'caption',
          description: 'description',
        },
      ],
      total: 1,
      loading: false,
    }));
    const tree = mount(
      <Provider store={store}>
        <Router>
          <Media permission={{ actions: ['create'] }} />
        </Router>
      </Provider>,
    )
    expect(tree).toMatchSnapshot();
  });
});
