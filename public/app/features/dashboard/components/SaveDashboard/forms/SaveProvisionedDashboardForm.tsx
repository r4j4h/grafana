import React, { useCallback, useState } from 'react';
import { css } from '@emotion/css';
import { saveAs } from 'file-saver';
import { Button, ClipboardButton, Modal, stylesFactory, TextArea, useTheme } from '@grafana/ui';
import { SaveDashboardFormProps } from '../types';
import { GrafanaTheme } from '@grafana/data';
import { useAppNotification } from 'app/core/copy/appNotification';

export const SaveProvisionedDashboardForm: React.FC<SaveDashboardFormProps> = ({ dashboard, onCancel }) => {
  const theme = useTheme();
  const notifyApp = useAppNotification();
  const [dashboardJSON, setDashboardJson] = useState(() => {
    const clone = dashboard.getSaveModelClone();
    delete clone.id;
    return JSON.stringify(clone, null, 2);
  });

  const saveToFile = useCallback(() => {
    const blob = new Blob([dashboardJSON], {
      type: 'application/json;charset=utf-8',
    });
    saveAs(blob, dashboard.title + '-' + new Date().getTime() + '.json');
  }, [dashboard.title, dashboardJSON]);

  const onCopyToClipboardSuccess = useCallback(() => {
    notifyApp.success('Dashboard JSON copied to clipboard');
  }, [notifyApp]);

  const styles = getStyles(theme);
  return (
    <>
      <div>
        <div>
          This dashboard cannot be saved from the Grafana UI because it has been provisioned from another source. Copy
          the JSON or save it to a file below, then you can update your dashboard in the provisioning source.
          <br />
          <i>
            See{' '}
            <a
              className="external-link"
              href="https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards"
              target="_blank"
              rel="noreferrer"
            >
              documentation
            </a>{' '}
            for more information about provisioning.
          </i>
          <br /> <br />
          <strong>File path: </strong> {dashboard.meta.provisionedExternalId}
        </div>
        <TextArea
          spellCheck={false}
          value={dashboardJSON}
          onChange={(e) => {
            setDashboardJson(e.currentTarget.value);
          }}
          className={styles.json}
        />
        <Modal.ButtonRow>
          <Button variant="secondary" onClick={onCancel} fill="outline">
            Cancel
          </Button>
          <ClipboardButton getText={() => dashboardJSON} onClipboardCopy={onCopyToClipboardSuccess}>
            Copy JSON to clipboard
          </ClipboardButton>
          <Button onClick={saveToFile}>Save JSON to file</Button>
        </Modal.ButtonRow>
      </div>
    </>
  );
};

const getStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    json: css`
      height: 400px;
      width: 100%;
      overflow: auto;
      resize: none;
      font-family: monospace;
    `,
  };
});
