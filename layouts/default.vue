<template>
  <div>
    <b-navbar toggleable="md" type="dark" variant="dark" fixed="top" sticky>
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-navbar-brand to="/list/general">ML News</b-navbar-brand>
      <b-collapse is-nav id="nav_collapse">
        <b-navbar-nav>
          <b-nav-item to="/list/general">General</b-nav-item>
          <b-nav-item to="/list/article">Article</b-nav-item>
          <b-nav-item to="/list/github">Github</b-nav-item>
          <b-nav-item to="/list/arxiv">arXiv</b-nav-item>
          <b-nav-item to="/list/slide">Slide</b-nav-item>
          <b-nav-item to="/search">Search</b-nav-item>
          <b-nav-item to="/recent-added-examples">recent-added-examples</b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <nuxt/>
  </div>
</template>
