import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";

// Views
import Home from "../views/Home.vue";
import Servers from "../views/Servers.vue";
import ServerDetail from "../views/ServerDetail.vue";
import Settings from "../views/Settings.vue";
import NotFound from "../views/NotFound.vue";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Home",
    component: Home,
    meta: { title: "Dashboard" },
  },
  {
    path: "/servers",
    name: "Servers",
    component: Servers,
    meta: { title: "Servers" },
  },
  {
    path: "/servers/:id",
    name: "ServerDetail",
    component: ServerDetail,
    meta: { title: "Server Details" },
  },
  {
    path: "/settings",
    name: "Settings",
    component: Settings,
    meta: { title: "Settings" },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: NotFound,
    meta: { title: "404 - Not Found" },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

// Update page title
router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title} | MCS Manager`;
  next();
});

export default router;
