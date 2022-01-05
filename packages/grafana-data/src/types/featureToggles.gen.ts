// NOTE: This file was auto generated.  DO NOT EDIT DIRECTLY!
// To change feature flags, edit:
//  pkg/setting/setting_feature_toggles_registry.go

/**
 * Describes available feature toggles in Grafana. These can be configured via
 * conf/custom.ini to enable features under development or not yet available in
 * stable version.
 *
 * @public
 */
export interface FeatureToggles {
  // [name: string]?: boolean; // support any string value

  recordedQueries?: boolean;
  trimDefaults?: boolean;
  database_metrics?: boolean;
  dashboardPreviews?: boolean;
  ['live-config']?: boolean;
  ['live-pipeline']?: boolean;
  ['live-service-web-worker']?: boolean;
  queryOverLive?: boolean;
  tempoSearch?: boolean;
  tempoServiceGraph?: boolean;
  fullRangeLogsVolume?: boolean;
  accesscontrol?: boolean;
  prometheus_azure_auth?: boolean;
  newNavigation?: boolean;
}