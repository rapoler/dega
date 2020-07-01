import axios from 'axios';
import {
  ADD_MEDIUM,
  ADD_MEDIA,
  ADD_MEDIA_REQUEST,
  SET_MEDIA_LOADING,
  RESET_MEDIA,
  MEDIA_API,
} from '../constants/media';
import { addErrors } from './notifications';

export const getMedia = (query) => {
  return (dispatch) => {
    dispatch(loadingMedia());
    return axios
      .get(MEDIA_API, {
        params: query,
      })
      .then((response) => {
        dispatch(addMediaList(response.data.nodes));
        dispatch(
          addMediaRequest({
            data: response.data.nodes.map((item) => item.id),
            query: query,
            total: response.data.total,
          }),
        );
        dispatch(stopMediaLoading());
      })
      .catch((error) => {
        console.log(error.message);
        dispatch(addErrors(error.message));
      });
  };
};

export const getMedium = (id) => {
  return (dispatch) => {
    dispatch(loadingMedia());
    return axios
      .get(MEDIA_API + '/' + id)
      .then((response) => {
        dispatch(getMediumByID(response.data));
        dispatch(stopMediaLoading());
      })
      .catch((error) => {
        dispatch(addErrors(error.message));
      });
  };
};

export const addMedium = (data) => {
  return (dispatch) => {
    dispatch(loadingMedia());
    return axios
      .post(MEDIA_API, data)
      .then(() => {
        dispatch(resetMedia());
      })
      .catch((error) => {
        dispatch(addErrors(error.message));
      });
  };
};

export const updateMedium = (data) => {
  return (dispatch) => {
    dispatch(loadingMedia());
    return axios
      .put(MEDIA_API + '/' + data.id, data)
      .then((response) => {
        dispatch(getMediumByID(response.data));
        dispatch(stopMediaLoading());
      })
      .catch((error) => {
        dispatch(addErrors(error.message));
      });
  };
};

export const deleteMedium = (id) => {
  return (dispatch) => {
    dispatch(loadingMedia());
    return axios
      .delete(MEDIA_API + '/' + id)
      .then(() => {
        dispatch(resetMedia());
      })
      .catch((error) => {
        dispatch(addErrors(error.message));
      });
  };
};

export const addMedia = (media) => {
  return (dispatch) => {
    dispatch(addMediaList(media));
  };
};

const loadingMedia = () => ({
  type: SET_MEDIA_LOADING,
  payload: true,
});

const stopMediaLoading = () => ({
  type: SET_MEDIA_LOADING,
  payload: false,
});

const getMediumByID = (data) => ({
  type: ADD_MEDIUM,
  payload: data,
});

export const addMediaList = (data) => ({
  type: ADD_MEDIA,
  payload: data,
});

const addMediaRequest = (data) => ({
  type: ADD_MEDIA_REQUEST,
  payload: data,
});

const resetMedia = () => ({
  type: RESET_MEDIA,
});
