import { createRouter, createWebHistory } from "vue-router";

import LoginView from "@/views/LoginView.vue";
import DashboardView from "@/views/DashboardView.vue";

import ProjectManageView from "@/views/project/ProjectManageView.vue";

import RequirementBookView from "@/views/requirement/RequirementBookView.vue";
import RequirementChangeView from "@/views/requirement/RequirementChangeView.vue";
import CustomerSuppliedView from "@/views/requirement/CustomerSuppliedView.vue";

import HardwareVersionView from "@/views/hardware/HardwareVersionView.vue";
import HardwareTestView from "@/views/hardware/HardwareTestView.vue";

import SoftwareVersionView from "@/views/software/SoftwareVersionView.vue";
import ProjectBranchView from "@/views/software/ProjectBranchView.vue";
// import ProjectConfigFileView from "@/views/software/ProjectConfigFileView.vue";

import TestCaseView from "@/views/test/TestCaseView.vue";
// import TestReportView from "@/views/test/TestReportView.vue";
import IssueCloseLoopView from "@/views/test/IssueCloseLoopView.vue";

import ProductionBurnRecordView from "@/views/production/ProductionBurnRecordView.vue";
import FactoryTestView from "@/views/production/FactoryTestView.vue";
import InventoryView from "@/views/production/InventoryView.vue";

import ShippingBatchView from "@/views/shipping/ShippingBatchView.vue";
// import ReceiptRecordView from "@/views/shipping/ReceiptRecordView.vue";
import OutInventoryView from "@/views/shipping/OutInventoryView.vue";

// import UpgradeRecordView from "@/views/aftersales/UpgradeRecordView.vue";
import RepairRecordView from "@/views/aftersales/RepairRecordView.vue";
import FaultAnalysisView from "@/views/aftersales/FaultAnalysisView.vue";

import ProjectProgressReportView from "@/views/report/ProjectProgressReportView.vue";
import VersionMatrixView from "@/views/report/VersionMatrixView.vue";
import IssueStatisticsView from "@/views/report/IssueStatisticsView.vue";
import { canAccessPage, getDefaultAccessiblePage } from "@/utils/permission";

const routes = [
  { path: "/login", component: LoginView },
  { path: "/", redirect: "/dashboard" },
  { path: "/dashboard", component: DashboardView },
  {
    path: '/project/progress-report',
    name: 'ProjectProgressReport',
    component: () => import('@/views/report/ProjectProgressReportView.vue')
  },
  {
    path: '/version/matrix',
    name: 'VersionMatrix',
    component: () => import('@/views/report/VersionMatrixView.vue')
  },
  { path: "/project/manage", component: ProjectManageView },

  { path: "/requirement/book", component: RequirementBookView },
  { path: "/requirement/change", component: RequirementChangeView },
  { path: "/requirement/customer-supplied", component: CustomerSuppliedView },
  { path: "/hardware/version", component: HardwareVersionView },
  { path: "/hardware/test", component: HardwareTestView },

  { path: "/software/version", component: SoftwareVersionView },
  { path: "/software/branch", component: ProjectBranchView },
  // { path: "/software/project-config", component: ProjectConfigFileView },

  { path: "/test/case", component: TestCaseView },
  // { path: "/test/report", component: TestReportView },
  { path: "/test/issue", component: IssueCloseLoopView },

  { path: "/production/burn", component: ProductionBurnRecordView },
  { path: "/production/factory-test", component: FactoryTestView },
  { path: "/production/inventory", component: InventoryView },
  {
    path: '/inventory/out',
    name: 'OutInventory',
    component: () => import('@/views/shipping/OutInventoryView.vue')
  },

  { path: "/shipping/batch", component: ShippingBatchView },
  // { path: "/shipping/receipt", component: ReceiptRecordView },
  { path: "/shipping/out", component: OutInventoryView },
  // { path: "/aftersales/upgrade", component: UpgradeRecordView },
  { path: "/aftersales/repair", component: RepairRecordView },
  { path: "/aftersales/fault-analysis", component: FaultAnalysisView },

  { path: "/report/project-progress", component: ProjectProgressReportView },
  { path: "/report/version-matrix", component: VersionMatrixView },
  { path: "/report/issue-statistics", component: IssueStatisticsView }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.path !== '/login' && !token) {
    next('/login')
    return
  }

  if (to.path === '/login' && token) {
    next(getDefaultAccessiblePage())
    return
  }

  if (token && !canAccessPage(to.path)) {
    next(getDefaultAccessiblePage())
    return
  }

  next()
})

export default router;
