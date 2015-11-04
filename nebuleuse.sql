-- phpMyAdmin SQL Dump
-- version 4.2.11
-- http://www.phpmyadmin.net
--
-- Client :  127.0.0.1
-- Généré le :  Mer 04 Novembre 2015 à 01:20
-- Version du serveur :  5.6.21
-- Version de PHP :  5.6.3

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- Base de données :  `nebuleuse`
--

-- --------------------------------------------------------

--
-- Structure de la table `neb_achievements`
--

CREATE TABLE IF NOT EXISTS `neb_achievements` (
`id` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `max` int(10) unsigned NOT NULL,
  `fullName` varchar(255) NOT NULL,
  `fullDesc` varchar(255) NOT NULL,
  `icon` varchar(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_achievements`
--

INSERT INTO `neb_achievements` (`id`, `name`, `max`, `fullName`, `fullDesc`, `icon`) VALUES
(1, 'test', 24, 'test', 'test', 'http://i.imgur.com/oyrwt3a.png');

-- --------------------------------------------------------

--
-- Structure de la table `neb_config`
--

CREATE TABLE IF NOT EXISTS `neb_config` (
  `name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_config`
--

INSERT INTO `neb_config` (`name`, `value`) VALUES
('gameName', ''),
('gameVersion', '1'),
('updaterVersion', '1'),
('sessionTimeout', '1800'),
('autoRegister', 'true'),
('defaultAvatar', 'http://i.imgur.com/oyrwt3a.png'),
('currentCommit', 'fbbf884cf75c703d8ead57dfaf5c0ecdb4ec37d1'),
('productionBranch', 'master'),
('gitRepositoryPath', './repo'),
('updateSystem', 'GitPatch');

-- --------------------------------------------------------

--
-- Structure de la table `neb_mirrors`
--

CREATE TABLE IF NOT EXISTS `neb_mirrors` (
  `id` smallint(5) unsigned NOT NULL,
  `Address` varchar(255) NOT NULL,
  `Name` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Structure de la table `neb_sessions`
--

CREATE TABLE IF NOT EXISTS `neb_sessions` (
  `userid` int(10) unsigned NOT NULL,
  `lastAlive` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sessionId` varchar(36) NOT NULL,
  `sessionStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Structure de la table `neb_stats_tables`
--

CREATE TABLE IF NOT EXISTS `neb_stats_tables` (
  `tableName` varchar(255) NOT NULL,
  `fields` text NOT NULL,
  `autoCount` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_stats_tables`
--

INSERT INTO `neb_stats_tables` (`tableName`, `fields`, `autoCount`) VALUES
('users', '[{"Name":"timeplayed","Type":"int","Size":11}]', 0),
('kills', '[{"Name":"userid","Type":"int","Size":11},{"Name":"x","Type":"int","Size":11},{"Name":"y","Type":"int","Size":11},{"Name":"z","Type":"int","Size":11},{"Name":"map","Type":"string","Size":255},{"Name":"weapon","Type":"string","Size":255}]', 1);

-- --------------------------------------------------------

--
-- Structure de la table `neb_updates`
--

CREATE TABLE IF NOT EXISTS `neb_updates` (
  `version` int(11) NOT NULL,
  `SemVer` varchar(16) NOT NULL,
  `log` text NOT NULL,
  `size` int(11) NOT NULL,
  `url` varchar(255) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `commit` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Structure de la table `neb_users`
--

CREATE TABLE IF NOT EXISTS `neb_users` (
`id` int(11) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `rank` tinyint(3) NOT NULL,
  `avatars` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_users`
--

INSERT INTO `neb_users` (`id`, `username`, `password`, `rank`, `avatars`, `hash`) VALUES
(1, 'test', 'q4F_1BnvOQERMAtuwNHoocjO6DiHvt15ol2krqZ60v-NW-tb0_IooASPuZq6iv1tjjT60JIIhA1MZvjTcGhDqA==', 2, '', 'O-z_gcTHzvgoM3ndhFnKVbM-tUcnGZDz_o6mhkWFiL0VTnvCvHFVOBYnvBp23pbz1ZIafCoH_JO51gXlVkmf8w==');

-- --------------------------------------------------------

--
-- Structure de la table `neb_users_achievements`
--

CREATE TABLE IF NOT EXISTS `neb_users_achievements` (
  `userid` int(10) unsigned NOT NULL,
  `achievementid` int(10) unsigned NOT NULL,
  `progress` int(10) unsigned NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_users_achievements`
--

INSERT INTO `neb_users_achievements` (`userid`, `achievementid`, `progress`) VALUES
(2, 1, 24);

-- --------------------------------------------------------

--
-- Structure de la table `neb_users_stats`
--

CREATE TABLE IF NOT EXISTS `neb_users_stats` (
  `userid` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `value` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Contenu de la table `neb_users_stats`
--

INSERT INTO `neb_users_stats` (`userid`, `name`, `value`) VALUES
(2, 'kills', 11);

-- --------------------------------------------------------

--
-- Structure de la table `neb_users_stats_kills`
--

CREATE TABLE IF NOT EXISTS `neb_users_stats_kills` (
  `userid` int(11) DEFAULT NULL,
  `x` int(11) DEFAULT NULL,
  `y` int(11) DEFAULT NULL,
  `z` int(11) DEFAULT NULL,
  `map` varchar(255) DEFAULT NULL,
  `weapon` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Index pour les tables exportées
--

--
-- Index pour la table `neb_achievements`
--
ALTER TABLE `neb_achievements`
 ADD PRIMARY KEY (`id`);

--
-- Index pour la table `neb_mirrors`
--
ALTER TABLE `neb_mirrors`
 ADD PRIMARY KEY (`id`);

--
-- Index pour la table `neb_sessions`
--
ALTER TABLE `neb_sessions`
 ADD PRIMARY KEY (`userid`);

--
-- Index pour la table `neb_updates`
--
ALTER TABLE `neb_updates`
 ADD PRIMARY KEY (`version`);

--
-- Index pour la table `neb_users`
--
ALTER TABLE `neb_users`
 ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT pour les tables exportées
--

--
-- AUTO_INCREMENT pour la table `neb_achievements`
--
ALTER TABLE `neb_achievements`
MODIFY `id` int(10) unsigned NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=2;
--
-- AUTO_INCREMENT pour la table `neb_users`
--
ALTER TABLE `neb_users`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=2;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
