import React, { useState } from 'react';
import { ConfigProvider, Button, Space, Tag, Table, Typography, Modal } from 'antd';
import {
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ExceptionOutlined,
  ClockCircleOutlined,
  CloseOutlined,
  FormOutlined,
} from '@ant-design/icons';
import { useDispatch } from 'react-redux';
import { deletePage } from '../../../actions/pages';
import { Link } from 'react-router-dom';
import QuickEdit from '../../../components/List/QuickEdit';

function PageList({ actions, format, status, data, filters, setFilters, fetchPages }) {
  const dispatch = useDispatch();
  const [id, setID] = useState(0);
  const [expandedRowKeys, setExpandedRowKeys] = useState([0]);
  const [modalOpen, setModalOpen] = useState(false);

  const getTagList = (tagids) => {
    return tagids.map((id) => (
      <Link to={`/pages?tag=${id}`}>
        <Tag>{data.tags[id].name}</Tag>
      </Link>
    ));
  };
  const getCategoryList = (catIds) => {
    return catIds.map((id) => (
      <Link to={`/pages?category=${id}`}>
        <Tag>{data.categories[id].name}</Tag>
      </Link>
    ));
  };
  // const getAuthorsList = (ids) => {
  //   return ids?.map((id) => (
  //     <Link>
  //       <Tag>{data.authors[id].display_name}</Tag>
  //     </Link>
  //   ));
  // };
  const columns = [
    {
      title: 'Title',
      dataIndex: 'title',
      key: 'title',
      width: 400,
      render: (_, item) => (
        <Link to={`/pages/${item.id}/edit`}>
          <Typography.Text style={{ fontSize: '1rem' }} strong>
            {item.title}
          </Typography.Text>
          {/*
          {item.published_date && (
            <p style={{ color: 'CaptionText' }}>
              Published on {dayjs(item.published_date).format('MMMM Do YYYY')}
            </p>
          )}
          <p style={{ color: 'CaptionText' }}>by {getAuthorsList(item.authors)}</p>
          */}
        </Link>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      width: 200,
      render: (status) => {
        return status === 'publish' ? (
          <Tag icon={<CheckCircleOutlined />} color="green">
            Published
          </Tag>
        ) : status === 'draft' ? (
          <Tag color="red" icon={<ExceptionOutlined />}>
            Draft
          </Tag>
        ) : status === 'ready' ? (
          <Tag color="gold" icon={<ClockCircleOutlined />}>
            Ready to Publish
          </Tag>
        ) : null;
      },
    },
    {
      title: 'Actions',
      dataIndex: 'actions',
      fixed: 'right',
      width: 240,
      render: (_, item, idx) => {
        const isOpen = item.id === expandedRowKeys[0];
        return (
          <ConfigProvider theme={{
            components: {
              Button: {
                controlHeight: 35,
                colorBorder: "#F2F2F2",
                colorPrimaryHover: "#F2F2F2"
              }
            }
          }}>
            <div style={{ display: 'flex', gap: "0.5rem" }}>
              <Link style={{ display: 'block' }} to={`/pages/${item.id}/edit`}>
                <Button
                  size="large"
                  icon={<EditOutlined style={{ color: "#858585" }} />}
                  disabled={!(actions.includes('admin') || actions.includes('update'))}
                />
              </Link>
              <Button
                size="large"
                disabled={!(actions.includes('admin') || actions.includes('update'))}
                onClick={() => {
                  isOpen ? setExpandedRowKeys([]) : setExpandedRowKeys([item.id]);
                  return setID(item.id);
                }}
                icon={isOpen ? <CloseOutlined style={{ color: "#858585" }} /> : <FormOutlined style={{ color: "#858585" }} />}
              />
              <Button
                size="large"
                onClick={() => { setModalOpen(true) }}
                icon={<DeleteOutlined style={{ color: "#858585" }} />}
                disabled={!(actions.includes('admin') || actions.includes('delete'))}
              />
              <Modal
                open={modalOpen}
                closable={false}
                centered
                width={400}
                className="delete-modal-container"
                style={{
                  borderRadius: '18px',
                }}
                onOk={() => {
                  () => dispatch(deletePage(item.id)).then(() => fetchPages())
                }}
                onCancel={() => {
                  setModalOpen(false);
                }}
              >
                <p>Are you sure you want to delete this Page?</p>
              </Modal>
              {/* <Button
                  icon={<EditOutlined />}
                  disabled={!(actions.includes('admin') || actions.includes('update'))}
                  style={{
                    margin: '0.5rem',
                    padding: '4px 22px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                /> */}
            </div>
          </ConfigProvider>
        );
      },
    },
  ];

  return (
    <Space direction="vertical">
      <Table
        dataSource={data.pages}
        loading={data.loading}
        columns={columns}
        rowKey={(record) => record.id}
        locale={{
          emptyText: '-',
        }}
        expandable={{
          expandIconColumnIndex: -1,
          expandedRowKeys,
          onExpand: (expanded, record) => {
            let keys = [];
            if (expanded) {
              keys.push(record.id);
            }

            setExpandedRowKeys(keys);
          },
          expandedRowRender: (item) => (
            <QuickEdit
              data={item}
              setID={setID}
              slug={format.slug}
              onQuickEditUpdate={() => setExpandedRowKeys([])}
            />
          ),
          expandIcon: () => { },
        }}
        pagination={{
          showTotal: (total, range) => `${range[0]}-${range[1]} of ${total} results`,
          total: data.total,
          current: filters.page,
          pageSize: filters.limit ? filters.limit : 10,
          onChange: (pageNumber, pageSize) =>
            setFilters({ ...filters, page: pageNumber, limit: pageSize }),
          pageSizeOptions: ['10', '15', '20'],
        }}
      />
    </Space>
  );
}

export default PageList;
