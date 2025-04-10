import { createRouter, createWebHistory } from 'vue-router';
import Workflow from '../components/Workflow.vue';
import WorkflowList from '../views/WorkflowList.vue';

const routes = [
  {
    path: '/',
    redirect: '/workflows'
  },
  {
    path: '/workflows',
    name: 'workflowList',
    component: WorkflowList
  },
  {
    path: '/workflow/create',
    name: 'createWorkflow',
    component: Workflow
  },
  {
    path: '/workflow/edit/:id',
    name: 'editWorkflow',
    component: Workflow,
    props: true
  },
  {
    path: '/workflow/statistics/:id',
    name: 'WorkflowStatistics',
    component: () => import('../views/WorkflowStatistics.vue')
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
